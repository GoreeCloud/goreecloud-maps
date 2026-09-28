package store

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const mapsTestRuntimeRole = "maps_test_runtime"

func TestMultiUserRLSIsolation(t *testing.T) {
	adminURL := strings.TrimSpace(os.Getenv("MAPS_TEST_DATABASE_URL"))
	if adminURL == "" {
		t.Skip("MAPS_TEST_DATABASE_URL is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	admin, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}
	defer admin.Close()

	if err := resetTestDatabase(ctx, admin); err != nil {
		t.Fatalf("reset test database: %v", err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if _, err := admin.Exec(cleanupCtx, "DROP SCHEMA IF EXISTS maps CASCADE"); err != nil {
			t.Logf("cleanup maps schema: %v", err)
		}
		if _, err := admin.Exec(cleanupCtx, "DROP ROLE IF EXISTS "+mapsTestRuntimeRole); err != nil {
			t.Logf("cleanup runtime role: %v", err)
		}
	}()

	if err := applyMigrations(ctx, admin, migrationDir(t)); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	if ownerStore, err := Open(ctx, adminURL); err == nil {
		ownerStore.Close()
		t.Fatal("database-owner/BYPASSRLS connection must be rejected")
	}

	runtimeURL, err := createRuntimeRole(ctx, admin, adminURL)
	if err != nil {
		t.Fatalf("create runtime role: %v", err)
	}

	dataStore, err := Open(ctx, runtimeURL)
	if err != nil {
		t.Fatalf("open non-owner runtime store: %v", err)
	}
	defer dataStore.Close()

	owner := resolveTestUser(t, ctx, dataStore, "maps-test-owner")
	editor := resolveTestUser(t, ctx, dataStore, "maps-test-editor")
	viewer := resolveTestUser(t, ctx, dataStore, "maps-test-viewer")
	stranger := resolveTestUser(t, ctx, dataStore, "maps-test-stranger")

	collection, err := dataStore.CreateCollection(ctx, owner.ID, CreateCollectionInput{
		Name:        "Chicago trip",
		Description: "Multi-user RLS integration fixture",
	})
	if err != nil {
		t.Fatalf("owner create collection: %v", err)
	}
	if collection.Revision != 1 || collection.Role != "owner" {
		t.Fatalf("unexpected new collection: %#v", collection)
	}

	mustCollectionCount(t, ctx, dataStore, owner.ID, 1)
	mustCollectionCount(t, ctx, dataStore, editor.ID, 0)
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 0)
	mustCollectionCount(t, ctx, dataStore, stranger.ID, 0)

	if err := dataStore.withUserTx(ctx, stranger.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.audit_events (actor_user_id, collection_id, event_type, event_data)
			VALUES ($1::uuid, $2::uuid, 'forged.audit', '{}'::jsonb)
		`, stranger.ID, collection.ID)
		return err
	}); err == nil {
		t.Fatal("stranger must not forge audit events for an inaccessible collection")
	}

	if _, err := dataStore.AddCollectionMember(ctx, owner.ID, collection.ID, editor.ID, "editor"); err != nil {
		t.Fatalf("owner add editor: %v", err)
	}
	if _, err := dataStore.AddCollectionMember(ctx, owner.ID, collection.ID, viewer.ID, "viewer"); err != nil {
		t.Fatalf("owner add viewer: %v", err)
	}
	if _, err := dataStore.AddCollectionMember(ctx, editor.ID, collection.ID, stranger.ID, "viewer"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("editor membership administration must be forbidden, got %v", err)
	}

	members, err := dataStore.ListCollectionMembers(ctx, viewer.ID, collection.ID)
	if err != nil {
		t.Fatalf("viewer list members: %v", err)
	}
	if len(members) != 3 {
		t.Fatalf("expected owner/editor/viewer membership view, got %#v", members)
	}

	mustCollectionCount(t, ctx, dataStore, editor.ID, 1)
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 1)
	mustCollectionCount(t, ctx, dataStore, stranger.ID, 0)

	editedCollection, err := dataStore.UpdateCollection(ctx, editor.ID, collection.ID, UpdateCollectionInput{
		Name:             stringPointer("Edited Chicago trip"),
		ExpectedRevision: collection.Revision,
	})
	if err != nil {
		t.Fatalf("editor update collection: %v", err)
	}
	if editedCollection.Revision != 2 || editedCollection.Name != "Edited Chicago trip" {
		t.Fatalf("unexpected edited collection: %#v", editedCollection)
	}
	if _, err := dataStore.UpdateCollection(ctx, viewer.ID, collection.ID, UpdateCollectionInput{
		Name:             stringPointer("Viewer edit must fail"),
		ExpectedRevision: editedCollection.Revision,
	}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("viewer collection update must be forbidden, got %v", err)
	}
	if _, err := dataStore.UpdateCollection(ctx, owner.ID, collection.ID, UpdateCollectionInput{
		Description:      stringPointer("stale write"),
		ExpectedRevision: 1,
	}); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale collection update must conflict, got %v", err)
	}

	item, err := dataStore.CreateCollectionItem(ctx, editor.ID, collection.ID, CreateCollectionItemInput{
		Provider:        "test",
		ProviderPlaceID: "union-station",
		Name:            "Union Station",
		Address:         "Chicago",
		Latitude:        41.881,
		Longitude:       -87.64,
		Note:            "Editor-created fixture",
		SortKey:         10,
	})
	if err != nil {
		t.Fatalf("editor create collection item: %v", err)
	}
	if item.Revision != 1 || item.CreatedByUserID != editor.ID {
		t.Fatalf("unexpected item: %#v", item)
	}
	if _, err := dataStore.CreateCollectionItem(ctx, viewer.ID, collection.ID, CreateCollectionItemInput{
		Provider:  "test",
		Name:      "Viewer write must fail",
		Latitude:  41.89,
		Longitude: -87.63,
	}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("viewer item creation must be forbidden, got %v", err)
	}

	ownerItems, err := dataStore.ListCollectionItems(ctx, owner.ID, collection.ID)
	if err != nil || len(ownerItems) != 1 {
		t.Fatalf("owner list collection items: items=%#v err=%v", ownerItems, err)
	}
	viewerItems, err := dataStore.ListCollectionItems(ctx, viewer.ID, collection.ID)
	if err != nil || len(viewerItems) != 1 {
		t.Fatalf("viewer list collection items: items=%#v err=%v", viewerItems, err)
	}
	if _, err := dataStore.ListCollectionItems(ctx, stranger.ID, collection.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("stranger must not enumerate collection items, got %v", err)
	}

	updatedItem, err := dataStore.UpdateCollectionItem(ctx, editor.ID, collection.ID, item.ID, UpdateCollectionItemInput{
		Note:             stringPointer("Updated note"),
		SortKey:          int64Pointer(20),
		ExpectedRevision: item.Revision,
	})
	if err != nil {
		t.Fatalf("editor update collection item: %v", err)
	}
	if updatedItem.Revision != 2 || updatedItem.Note != "Updated note" || updatedItem.SortKey != 20 {
		t.Fatalf("unexpected updated item: %#v", updatedItem)
	}
	if _, err := dataStore.UpdateCollectionItem(ctx, editor.ID, collection.ID, item.ID, UpdateCollectionItemInput{
		Note:             stringPointer("stale item update"),
		ExpectedRevision: 1,
	}); !errors.Is(err, ErrConflict) {
		t.Fatalf("stale item update must conflict, got %v", err)
	}
	if _, err := dataStore.UpdateCollectionItem(ctx, viewer.ID, collection.ID, item.ID, UpdateCollectionItemInput{
		Note:             stringPointer("viewer edit"),
		ExpectedRevision: updatedItem.Revision,
	}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("viewer item update must be forbidden, got %v", err)
	}

	if _, err := dataStore.UpdateCollectionMemberRole(ctx, owner.ID, collection.ID, editor.ID, "viewer"); err != nil {
		t.Fatalf("owner demote editor: %v", err)
	}
	if _, err := dataStore.UpdateCollectionItem(ctx, editor.ID, collection.ID, item.ID, UpdateCollectionItemInput{
		Note:             stringPointer("demoted editor write"),
		ExpectedRevision: updatedItem.Revision,
	}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("demoted editor must not update item, got %v", err)
	}
	if _, err := dataStore.UpdateCollectionMemberRole(ctx, owner.ID, collection.ID, editor.ID, "editor"); err != nil {
		t.Fatalf("owner restore editor: %v", err)
	}

	if err := dataStore.withUserTx(ctx, owner.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.saved_places (
				owner_user_id, provider, provider_place_id, name, position
			)
			VALUES (
				$1, 'test', 'private-place', 'Private place',
				ST_SetSRID(ST_MakePoint(-87.62, 41.90), 4326)::geography
			)
		`, owner.ID)
		return err
	}); err != nil {
		t.Fatalf("owner create saved place: %v", err)
	}
	mustSavedPlaceCount(t, ctx, dataStore, owner.ID, 1)
	mustSavedPlaceCount(t, ctx, dataStore, stranger.ID, 0)

	if err := dataStore.RemoveCollectionMember(ctx, viewer.ID, collection.ID, viewer.ID); err != nil {
		t.Fatalf("viewer self-remove membership: %v", err)
	}
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 0)
	if _, err := dataStore.ListCollectionMembers(ctx, viewer.ID, collection.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("removed viewer must lose collection access, got %v", err)
	}

	if err := dataStore.DeleteCollectionItem(ctx, editor.ID, collection.ID, item.ID, updatedItem.Revision); err != nil {
		t.Fatalf("editor delete collection item: %v", err)
	}
	itemsAfterDelete, err := dataStore.ListCollectionItems(ctx, owner.ID, collection.ID)
	if err != nil || len(itemsAfterDelete) != 0 {
		t.Fatalf("expected deleted item to disappear: items=%#v err=%v", itemsAfterDelete, err)
	}

	if err := dataStore.withUserTx(ctx, owner.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE maps.collections SET owner_user_id = $2 WHERE id = $1`, collection.ID, editor.ID)
		return err
	}); err == nil {
		t.Fatal("collection ownership mutation must be rejected")
	}

	var auditCount int
	if err := dataStore.withUserTx(ctx, owner.ID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT count(*) FROM maps.audit_events WHERE collection_id = $1`, collection.ID).Scan(&auditCount)
	}); err != nil {
		t.Fatalf("owner read audit events: %v", err)
	}
	if auditCount < 8 {
		t.Fatalf("expected collaboration audit events, got %d", auditCount)
	}
}

func resetTestDatabase(ctx context.Context, admin *pgxpool.Pool) error {
	if _, err := admin.Exec(ctx, "DROP SCHEMA IF EXISTS maps CASCADE"); err != nil {
		return err
	}
	_, err := admin.Exec(ctx, "DROP ROLE IF EXISTS "+mapsTestRuntimeRole)
	return err
}

func createRuntimeRole(ctx context.Context, admin *pgxpool.Pool, adminURL string) (string, error) {
	statements := []string{
		"CREATE ROLE " + mapsTestRuntimeRole + " LOGIN PASSWORD 'maps-test-runtime' NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT NOBYPASSRLS",
		"GRANT USAGE ON SCHEMA maps TO " + mapsTestRuntimeRole,
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA maps TO " + mapsTestRuntimeRole,
		"GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA maps TO " + mapsTestRuntimeRole,
		"GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA maps TO " + mapsTestRuntimeRole,
	}
	for _, statement := range statements {
		if _, err := admin.Exec(ctx, statement); err != nil {
			return "", err
		}
	}

	parsed, err := url.Parse(adminURL)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return "", fmt.Errorf("MAPS_TEST_DATABASE_URL must be a PostgreSQL URL")
	}
	parsed.User = url.UserPassword(mapsTestRuntimeRole, "maps-test-runtime")
	return parsed.String(), nil
}

func resolveTestUser(t *testing.T, ctx context.Context, dataStore *Store, subject string) User {
	t.Helper()
	user, err := dataStore.ResolveUser(ctx, subject)
	if err != nil {
		t.Fatalf("resolve %s: %v", subject, err)
	}
	return user
}

func mustCollectionCount(t *testing.T, ctx context.Context, dataStore *Store, userID string, expected int) {
	t.Helper()
	collections, err := dataStore.ListCollections(ctx, userID)
	if err != nil {
		t.Fatalf("list collections for %s: %v", userID, err)
	}
	if len(collections) != expected {
		t.Fatalf("expected %d visible collections for %s, got %d", expected, userID, len(collections))
	}
}

func mustSavedPlaceCount(t *testing.T, ctx context.Context, dataStore *Store, userID string, expected int) {
	t.Helper()
	var count int
	if err := dataStore.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT count(*) FROM maps.saved_places`).Scan(&count)
	}); err != nil {
		t.Fatalf("count saved places for %s: %v", userID, err)
	}
	if count != expected {
		t.Fatalf("expected %d visible saved places for %s, got %d", expected, userID, count)
	}
}

func migrationDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../..", "db", "migrations"))
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		migration, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		if err := execMultiStatement(ctx, pool, string(migration)); err != nil {
			return fmt.Errorf("apply %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func execMultiStatement(ctx context.Context, pool *pgxpool.Pool, sql string) error {
	connection, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	results, err := connection.Conn().PgConn().Exec(ctx, sql).ReadAll()
	if err != nil {
		return err
	}
	for _, result := range results {
		if result.Err != nil {
			return result.Err
		}
	}
	return nil
}

func stringPointer(value string) *string {
	return &value
}

func int64Pointer(value int64) *int64 {
	return &value
}
