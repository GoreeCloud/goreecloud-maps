package store

import (
	"context"
	"fmt"
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

	migration, err := os.ReadFile(migrationPath(t))
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if err := execMultiStatement(ctx, admin, string(migration)); err != nil {
		t.Fatalf("apply migration: %v", err)
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
		Description: "Two-user RLS integration fixture",
	})
	if err != nil {
		t.Fatalf("owner create collection: %v", err)
	}

	mustCollectionCount(t, ctx, dataStore, owner.ID, 1)
	mustCollectionCount(t, ctx, dataStore, editor.ID, 0)
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 0)
	mustCollectionCount(t, ctx, dataStore, stranger.ID, 0)

	if err := dataStore.withUserTx(ctx, owner.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.collection_members (collection_id, user_id, role, invited_by_user_id)
			VALUES ($1, $2, 'editor', $4), ($1, $3, 'viewer', $4)
		`, collection.ID, editor.ID, viewer.ID, owner.ID)
		return err
	}); err != nil {
		t.Fatalf("owner add members: %v", err)
	}

	mustCollectionCount(t, ctx, dataStore, editor.ID, 1)
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 1)
	mustCollectionCount(t, ctx, dataStore, stranger.ID, 0)

	if err := dataStore.withUserTx(ctx, editor.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.collection_items (
				collection_id, created_by_user_id, provider, provider_place_id,
				name, address, position, note, sort_key
			)
			VALUES (
				$1, $2, 'test', 'union-station', 'Union Station', 'Chicago',
				ST_SetSRID(ST_MakePoint(-87.6400, 41.8810), 4326)::geography,
				'Editor-created fixture', 10
			)
		`, collection.ID, editor.ID)
		return err
	}); err != nil {
		t.Fatalf("editor add collection item: %v", err)
	}

	mustVisibleItemCount(t, ctx, dataStore, owner.ID, collection.ID, 1)
	mustVisibleItemCount(t, ctx, dataStore, editor.ID, collection.ID, 1)
	mustVisibleItemCount(t, ctx, dataStore, viewer.ID, collection.ID, 1)
	mustVisibleItemCount(t, ctx, dataStore, stranger.ID, collection.ID, 0)

	if err := dataStore.withUserTx(ctx, editor.ID, func(tx pgx.Tx) error {
		command, err := tx.Exec(ctx, `UPDATE maps.collections SET name = 'Edited Chicago trip' WHERE id = $1`, collection.ID)
		if err != nil {
			return err
		}
		if command.RowsAffected() != 1 {
			return fmt.Errorf("expected editor update to affect 1 row, affected %d", command.RowsAffected())
		}
		return nil
	}); err != nil {
		t.Fatalf("editor update collection: %v", err)
	}

	var viewerUpdateRows int64
	if err := dataStore.withUserTx(ctx, viewer.ID, func(tx pgx.Tx) error {
		command, err := tx.Exec(ctx, `UPDATE maps.collections SET name = 'Viewer edit must fail' WHERE id = $1`, collection.ID)
		if err != nil {
			return err
		}
		viewerUpdateRows = command.RowsAffected()
		return nil
	}); err != nil {
		t.Fatalf("viewer update authorization check: %v", err)
	}
	if viewerUpdateRows != 0 {
		t.Fatalf("viewer must not update collection; affected %d rows", viewerUpdateRows)
	}
	mustCollectionName(t, ctx, dataStore, owner.ID, collection.ID, "Edited Chicago trip")

	if err := dataStore.withUserTx(ctx, viewer.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.collection_items (
				collection_id, created_by_user_id, provider, name, position
			)
			VALUES (
				$1, $2, 'test', 'Viewer write must fail',
				ST_SetSRID(ST_MakePoint(-87.63, 41.89), 4326)::geography
			)
		`, collection.ID, viewer.ID)
		return err
	}); err == nil {
		t.Fatal("viewer must not insert collection items")
	}

	if err := dataStore.withUserTx(ctx, editor.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO maps.collection_members (collection_id, user_id, role, invited_by_user_id)
			VALUES ($1, $2, 'viewer', $3)
		`, collection.ID, stranger.ID, editor.ID)
		return err
	}); err == nil {
		t.Fatal("editor must not manage collection membership")
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

	if err := dataStore.withUserTx(ctx, viewer.ID, func(tx pgx.Tx) error {
		command, err := tx.Exec(ctx, `DELETE FROM maps.collection_members WHERE collection_id = $1 AND user_id = $2`, collection.ID, viewer.ID)
		if err != nil {
			return err
		}
		if command.RowsAffected() != 1 {
			return fmt.Errorf("expected viewer self-removal to affect 1 row, affected %d", command.RowsAffected())
		}
		return nil
	}); err != nil {
		t.Fatalf("viewer self-remove membership: %v", err)
	}
	mustCollectionCount(t, ctx, dataStore, viewer.ID, 0)

	if err := dataStore.withUserTx(ctx, owner.ID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE maps.collections SET owner_user_id = $2 WHERE id = $1`, collection.ID, editor.ID)
		return err
	}); err == nil {
		t.Fatal("collection ownership mutation must be rejected")
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

	config, err := pgxpool.ParseConfig(adminURL)
	if err != nil {
		return "", err
	}
	config.ConnConfig.User = mapsTestRuntimeRole
	config.ConnConfig.Password = "maps-test-runtime"
	return config.ConnString(), nil
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

func mustCollectionName(t *testing.T, ctx context.Context, dataStore *Store, userID, collectionID, expected string) {
	t.Helper()
	var name string
	if err := dataStore.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT name FROM maps.collections WHERE id = $1`, collectionID).Scan(&name)
	}); err != nil {
		t.Fatalf("read collection name for %s: %v", userID, err)
	}
	if name != expected {
		t.Fatalf("expected collection name %q, got %q", expected, name)
	}
}

func mustVisibleItemCount(t *testing.T, ctx context.Context, dataStore *Store, userID, collectionID string, expected int) {
	t.Helper()
	var count int
	if err := dataStore.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT count(*) FROM maps.collection_items WHERE collection_id = $1`, collectionID).Scan(&count)
	}); err != nil {
		t.Fatalf("count collection items for %s: %v", userID, err)
	}
	if count != expected {
		t.Fatalf("expected %d visible items for %s, got %d", expected, userID, count)
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

func migrationPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve integration test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../..", "db", "migrations", "0001_multi_user_foundation.sql"))
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
