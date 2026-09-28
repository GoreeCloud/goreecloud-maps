package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrForbidden    = errors.New("operation forbidden")
	ErrConflict     = errors.New("revision conflict")
	ErrInvalidInput = errors.New("invalid input")
)

type Store struct {
	pool *pgxpool.Pool
}

type User struct {
	ID              string    `json:"id"`
	IdentitySubject string    `json:"-"`
	DisplayName     string    `json:"displayName,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SharingMode string    `json:"sharingMode"`
	Role        string    `json:"role"`
	Revision    int64     `json:"revision"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateCollectionInput struct {
	Name        string
	Description string
}

type UpdateCollectionInput struct {
	Name             *string
	Description      *string
	ExpectedRevision int64
}

type CollectionMember struct {
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName,omitempty"`
	Role        string    `json:"role"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CollectionItem struct {
	ID              string    `json:"id"`
	CollectionID    string    `json:"collectionId"`
	CreatedByUserID string    `json:"createdByUserId"`
	Provider        string    `json:"provider"`
	ProviderPlaceID string    `json:"providerPlaceId,omitempty"`
	Name            string    `json:"name"`
	Address         string    `json:"address,omitempty"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Note            string    `json:"note,omitempty"`
	SortKey         int64     `json:"sortKey"`
	Revision        int64     `json:"revision"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateCollectionItemInput struct {
	Provider        string
	ProviderPlaceID string
	Name            string
	Address         string
	Latitude        float64
	Longitude       float64
	Note            string
	SortKey         int64
}

type UpdateCollectionItemInput struct {
	Name             *string
	Address          *string
	Latitude         *float64
	Longitude        *float64
	Note             *string
	SortKey          *int64
	ExpectedRevision int64
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	var isTableOwner bool
	var bypassRLS bool
	if err := pool.QueryRow(ctx, `
		SELECT
			pg_has_role(current_user, c.relowner, 'MEMBER'),
			COALESCE((SELECT rolbypassrls FROM pg_roles WHERE rolname = current_user), false)
		FROM pg_class AS c
		JOIN pg_namespace AS n ON n.oid = c.relnamespace
		WHERE n.nspname = 'maps' AND c.relname = 'collections'
	`).Scan(&isTableOwner, &bypassRLS); err != nil {
		pool.Close()
		return nil, err
	}
	if isTableOwner || bypassRLS {
		pool.Close()
		return nil, errors.New("MAPS_DATABASE_URL must use a non-owner PostgreSQL role without BYPASSRLS")
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) ResolveUser(ctx context.Context, identitySubject string) (User, error) {
	identitySubject = strings.TrimSpace(identitySubject)
	if identitySubject == "" {
		return User{}, errors.New("identity subject is required")
	}

	var user User
	err := s.pool.QueryRow(ctx, `
		INSERT INTO maps.users (identity_subject)
		VALUES ($1)
		ON CONFLICT (identity_subject)
		DO UPDATE SET updated_at = maps.users.updated_at
		RETURNING id::text, identity_subject, COALESCE(display_name, ''), created_at
	`, identitySubject).Scan(&user.ID, &user.IdentitySubject, &user.DisplayName, &user.CreatedAt)
	return user, err
}

func (s *Store) ListCollections(ctx context.Context, userID string) ([]Collection, error) {
	collections := make([]Collection, 0)
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT
				id::text,
				name,
				COALESCE(description, ''),
				sharing_mode,
				maps.collection_access_role(id),
				revision,
				updated_at
			FROM maps.collections
			ORDER BY updated_at DESC, id
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var collection Collection
			if err := rows.Scan(
				&collection.ID,
				&collection.Name,
				&collection.Description,
				&collection.SharingMode,
				&collection.Role,
				&collection.Revision,
				&collection.UpdatedAt,
			); err != nil {
				return err
			}
			collections = append(collections, collection)
		}
		return rows.Err()
	})
	return collections, err
}

func (s *Store) CreateCollection(ctx context.Context, userID string, input CreateCollectionInput) (Collection, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	if err := validateCollectionText(input.Name, input.Description); err != nil {
		return Collection{}, err
	}

	var collection Collection
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, `
			INSERT INTO maps.collections (owner_user_id, name, description)
			VALUES ($1::uuid, $2, NULLIF($3, ''))
			RETURNING
				id::text,
				name,
				COALESCE(description, ''),
				sharing_mode,
				'owner',
				revision,
				updated_at
		`, userID, input.Name, input.Description).Scan(
			&collection.ID,
			&collection.Name,
			&collection.Description,
			&collection.SharingMode,
			&collection.Role,
			&collection.Revision,
			&collection.UpdatedAt,
		)
	})
	return collection, err
}

func (s *Store) UpdateCollection(ctx context.Context, userID, collectionID string, input UpdateCollectionInput) (Collection, error) {
	if input.ExpectedRevision < 1 || (input.Name == nil && input.Description == nil) {
		return Collection{}, ErrInvalidInput
	}

	var name any
	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if err := validateCollectionText(trimmed, ""); err != nil {
			return Collection{}, err
		}
		name = trimmed
	}
	var description any
	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		if len([]rune(trimmed)) > 4000 {
			return Collection{}, ErrInvalidInput
		}
		description = trimmed
	}

	var collection Collection
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}
		if role != "owner" && role != "editor" {
			return ErrForbidden
		}

		err = tx.QueryRow(ctx, `
			UPDATE maps.collections
			SET
				name = COALESCE($2::text, name),
				description = CASE WHEN $3::text IS NULL THEN description ELSE NULLIF($3::text, '') END,
				revision = revision + 1,
				updated_at = now()
			WHERE id = $1::uuid AND revision = $4
			RETURNING
				id::text,
				name,
				COALESCE(description, ''),
				sharing_mode,
				maps.collection_access_role(id),
				revision,
				updated_at
		`, collectionID, name, description, input.ExpectedRevision).Scan(
			&collection.ID,
			&collection.Name,
			&collection.Description,
			&collection.SharingMode,
			&collection.Role,
			&collection.Revision,
			&collection.UpdatedAt,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrConflict
		}
		if err != nil {
			return err
		}
		return insertAuditEvent(ctx, tx, userID, collectionID, "collection.updated", map[string]any{"revision": collection.Revision})
	})
	return collection, err
}

func (s *Store) ListCollectionMembers(ctx context.Context, userID, collectionID string) ([]CollectionMember, error) {
	members := make([]CollectionMember, 0)
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}

		rows, err := tx.Query(ctx, `
			SELECT c.owner_user_id::text, COALESCE(owner_user.display_name, ''), 'owner', c.updated_at
			FROM maps.collections AS c
			JOIN maps.users AS owner_user ON owner_user.id = c.owner_user_id
			WHERE c.id = $1::uuid
			UNION ALL
			SELECT cm.user_id::text, COALESCE(member_user.display_name, ''), cm.role, cm.updated_at
			FROM maps.collection_members AS cm
			JOIN maps.users AS member_user ON member_user.id = cm.user_id
			WHERE cm.collection_id = $1::uuid
			ORDER BY 3, 1
		`, collectionID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var member CollectionMember
			if err := rows.Scan(&member.UserID, &member.DisplayName, &member.Role, &member.UpdatedAt); err != nil {
				return err
			}
			members = append(members, member)
		}
		return rows.Err()
	})
	return members, err
}

func (s *Store) AddCollectionMember(ctx context.Context, actorUserID, collectionID, memberUserID, role string) (CollectionMember, error) {
	role = normalizeMemberRole(role)
	if role == "" || memberUserID == "" || memberUserID == actorUserID {
		return CollectionMember{}, ErrInvalidInput
	}

	var member CollectionMember
	err := s.withUserTx(ctx, actorUserID, func(tx pgx.Tx) error {
		actorRole, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if actorRole == "" {
			return ErrNotFound
		}
		if actorRole != "owner" {
			return ErrForbidden
		}

		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM maps.users WHERE id = $1::uuid)`, memberUserID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}

		err = tx.QueryRow(ctx, `
			INSERT INTO maps.collection_members (collection_id, user_id, role, invited_by_user_id)
			VALUES ($1::uuid, $2::uuid, $3, $4::uuid)
			ON CONFLICT (collection_id, user_id) DO NOTHING
			RETURNING user_id::text, '', role, updated_at
		`, collectionID, memberUserID, role, actorUserID).Scan(&member.UserID, &member.DisplayName, &member.Role, &member.UpdatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrConflict
		}
		if err != nil {
			return err
		}
		return insertAuditEvent(ctx, tx, actorUserID, collectionID, "collection.member_added", map[string]any{
			"member_user_id": memberUserID,
			"role":           role,
		})
	})
	return member, err
}

func (s *Store) UpdateCollectionMemberRole(ctx context.Context, actorUserID, collectionID, memberUserID, role string) (CollectionMember, error) {
	role = normalizeMemberRole(role)
	if role == "" || memberUserID == "" || memberUserID == actorUserID {
		return CollectionMember{}, ErrInvalidInput
	}

	var member CollectionMember
	err := s.withUserTx(ctx, actorUserID, func(tx pgx.Tx) error {
		actorRole, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if actorRole == "" {
			return ErrNotFound
		}
		if actorRole != "owner" {
			return ErrForbidden
		}

		err = tx.QueryRow(ctx, `
			UPDATE maps.collection_members
			SET role = $3, updated_at = now()
			WHERE collection_id = $1::uuid AND user_id = $2::uuid
			RETURNING user_id::text, '', role, updated_at
		`, collectionID, memberUserID, role).Scan(&member.UserID, &member.DisplayName, &member.Role, &member.UpdatedAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return insertAuditEvent(ctx, tx, actorUserID, collectionID, "collection.member_role_changed", map[string]any{
			"member_user_id": memberUserID,
			"role":           role,
		})
	})
	return member, err
}

func (s *Store) RemoveCollectionMember(ctx context.Context, actorUserID, collectionID, memberUserID string) error {
	if memberUserID == "" {
		return ErrInvalidInput
	}
	return s.withUserTx(ctx, actorUserID, func(tx pgx.Tx) error {
		actorRole, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if actorRole == "" {
			return ErrNotFound
		}
		if actorRole != "owner" && actorUserID != memberUserID {
			return ErrForbidden
		}

		var memberRole string
		if err := tx.QueryRow(ctx, `
			SELECT role
			FROM maps.collection_members
			WHERE collection_id = $1::uuid AND user_id = $2::uuid
		`, collectionID, memberUserID).Scan(&memberRole); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		if err := insertAuditEvent(ctx, tx, actorUserID, collectionID, "collection.member_removed", map[string]any{
			"member_user_id": memberUserID,
			"role":           memberRole,
		}); err != nil {
			return err
		}
		command, err := tx.Exec(ctx, `DELETE FROM maps.collection_members WHERE collection_id = $1::uuid AND user_id = $2::uuid`, collectionID, memberUserID)
		if err != nil {
			return err
		}
		if command.RowsAffected() != 1 {
			return ErrNotFound
		}
		return nil
	})
}

func (s *Store) ListCollectionItems(ctx context.Context, userID, collectionID string) ([]CollectionItem, error) {
	items := make([]CollectionItem, 0)
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}

		rows, err := tx.Query(ctx, `
			SELECT
				id::text,
				collection_id::text,
				created_by_user_id::text,
				provider,
				COALESCE(provider_place_id, ''),
				name,
				COALESCE(address, ''),
				ST_Y(position::geometry),
				ST_X(position::geometry),
				COALESCE(note, ''),
				sort_key,
				revision,
				created_at,
				updated_at
			FROM maps.collection_items
			WHERE collection_id = $1::uuid
			ORDER BY sort_key, created_at, id
		`, collectionID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item CollectionItem
			if err := scanCollectionItem(rows, &item); err != nil {
				return err
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	return items, err
}

func (s *Store) CreateCollectionItem(ctx context.Context, userID, collectionID string, input CreateCollectionItemInput) (CollectionItem, error) {
	input.Provider = strings.TrimSpace(input.Provider)
	input.ProviderPlaceID = strings.TrimSpace(input.ProviderPlaceID)
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Note = strings.TrimSpace(input.Note)
	if err := validateCollectionItem(input.Provider, input.ProviderPlaceID, input.Name, input.Address, input.Note, input.Latitude, input.Longitude); err != nil {
		return CollectionItem{}, err
	}

	var item CollectionItem
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}
		if role != "owner" && role != "editor" {
			return ErrForbidden
		}

		err = tx.QueryRow(ctx, `
			INSERT INTO maps.collection_items (
				collection_id, created_by_user_id, provider, provider_place_id,
				name, address, position, note, sort_key
			)
			VALUES (
				$1::uuid, $2::uuid, $3, NULLIF($4, ''),
				$5, NULLIF($6, ''),
				ST_SetSRID(ST_MakePoint($8, $7), 4326)::geography,
				NULLIF($9, ''), $10
			)
			RETURNING
				id::text,
				collection_id::text,
				created_by_user_id::text,
				provider,
				COALESCE(provider_place_id, ''),
				name,
				COALESCE(address, ''),
				ST_Y(position::geometry),
				ST_X(position::geometry),
				COALESCE(note, ''),
				sort_key,
				revision,
				created_at,
				updated_at
		`, collectionID, userID, input.Provider, input.ProviderPlaceID, input.Name, input.Address, input.Latitude, input.Longitude, input.Note, input.SortKey).Scan(
			&item.ID,
			&item.CollectionID,
			&item.CreatedByUserID,
			&item.Provider,
			&item.ProviderPlaceID,
			&item.Name,
			&item.Address,
			&item.Latitude,
			&item.Longitude,
			&item.Note,
			&item.SortKey,
			&item.Revision,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return err
		}
		return insertAuditEvent(ctx, tx, userID, collectionID, "collection.item_added", map[string]any{"item_id": item.ID})
	})
	return item, err
}

func (s *Store) UpdateCollectionItem(ctx context.Context, userID, collectionID, itemID string, input UpdateCollectionItemInput) (CollectionItem, error) {
	if input.ExpectedRevision < 1 || (input.Name == nil && input.Address == nil && input.Latitude == nil && input.Longitude == nil && input.Note == nil && input.SortKey == nil) {
		return CollectionItem{}, ErrInvalidInput
	}
	if (input.Latitude == nil) != (input.Longitude == nil) {
		return CollectionItem{}, ErrInvalidInput
	}

	var name any
	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if trimmed == "" || len([]rune(trimmed)) > 240 {
			return CollectionItem{}, ErrInvalidInput
		}
		name = trimmed
	}
	var address any
	if input.Address != nil {
		trimmed := strings.TrimSpace(*input.Address)
		if len([]rune(trimmed)) > 1000 {
			return CollectionItem{}, ErrInvalidInput
		}
		address = trimmed
	}
	var note any
	if input.Note != nil {
		trimmed := strings.TrimSpace(*input.Note)
		if len([]rune(trimmed)) > 4000 {
			return CollectionItem{}, ErrInvalidInput
		}
		note = trimmed
	}
	if input.Latitude != nil && !validCoordinates(*input.Latitude, *input.Longitude) {
		return CollectionItem{}, ErrInvalidInput
	}

	var item CollectionItem
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}
		if role != "owner" && role != "editor" {
			return ErrForbidden
		}

		var currentRevision int64
		if err := tx.QueryRow(ctx, `SELECT revision FROM maps.collection_items WHERE collection_id = $1::uuid AND id = $2::uuid`, collectionID, itemID).Scan(&currentRevision); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}
		if currentRevision != input.ExpectedRevision {
			return ErrConflict
		}

		err = tx.QueryRow(ctx, `
			UPDATE maps.collection_items
			SET
				name = COALESCE($3::text, name),
				address = CASE WHEN $4::text IS NULL THEN address ELSE NULLIF($4::text, '') END,
				position = CASE
					WHEN $5::double precision IS NULL THEN position
					ELSE ST_SetSRID(ST_MakePoint($6::double precision, $5::double precision), 4326)::geography
				END,
				note = CASE WHEN $7::text IS NULL THEN note ELSE NULLIF($7::text, '') END,
				sort_key = COALESCE($8::bigint, sort_key),
				revision = revision + 1,
				updated_at = now()
			WHERE collection_id = $1::uuid AND id = $2::uuid AND revision = $9
			RETURNING
				id::text,
				collection_id::text,
				created_by_user_id::text,
				provider,
				COALESCE(provider_place_id, ''),
				name,
				COALESCE(address, ''),
				ST_Y(position::geometry),
				ST_X(position::geometry),
				COALESCE(note, ''),
				sort_key,
				revision,
				created_at,
				updated_at
		`, collectionID, itemID, name, address, input.Latitude, input.Longitude, note, input.SortKey, input.ExpectedRevision).Scan(
			&item.ID,
			&item.CollectionID,
			&item.CreatedByUserID,
			&item.Provider,
			&item.ProviderPlaceID,
			&item.Name,
			&item.Address,
			&item.Latitude,
			&item.Longitude,
			&item.Note,
			&item.SortKey,
			&item.Revision,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrConflict
		}
		if err != nil {
			return err
		}
		return insertAuditEvent(ctx, tx, userID, collectionID, "collection.item_updated", map[string]any{
			"item_id":  item.ID,
			"revision": item.Revision,
		})
	})
	return item, err
}

func (s *Store) DeleteCollectionItem(ctx context.Context, userID, collectionID, itemID string, expectedRevision int64) error {
	if expectedRevision < 1 {
		return ErrInvalidInput
	}
	return s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		role, err := collectionRole(ctx, tx, collectionID)
		if err != nil {
			return err
		}
		if role == "" {
			return ErrNotFound
		}
		if role != "owner" && role != "editor" {
			return ErrForbidden
		}

		var currentRevision int64
		if err := tx.QueryRow(ctx, `SELECT revision FROM maps.collection_items WHERE collection_id = $1::uuid AND id = $2::uuid`, collectionID, itemID).Scan(&currentRevision); errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		} else if err != nil {
			return err
		}
		if currentRevision != expectedRevision {
			return ErrConflict
		}
		if err := insertAuditEvent(ctx, tx, userID, collectionID, "collection.item_deleted", map[string]any{"item_id": itemID}); err != nil {
			return err
		}
		command, err := tx.Exec(ctx, `DELETE FROM maps.collection_items WHERE collection_id = $1::uuid AND id = $2::uuid AND revision = $3`, collectionID, itemID, expectedRevision)
		if err != nil {
			return err
		}
		if command.RowsAffected() != 1 {
			return ErrConflict
		}
		return nil
	})
}

func (s *Store) withUserTx(ctx context.Context, userID string, fn func(pgx.Tx) error) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	if _, err := tx.Exec(ctx, `SELECT set_config('goreecloud.maps_user_id', $1, true)`, userID); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func collectionRole(ctx context.Context, tx pgx.Tx, collectionID string) (string, error) {
	var role *string
	if err := tx.QueryRow(ctx, `SELECT maps.collection_access_role($1::uuid)`, collectionID).Scan(&role); err != nil {
		return "", err
	}
	if role == nil {
		return "", nil
	}
	return *role, nil
}

func insertAuditEvent(ctx context.Context, tx pgx.Tx, actorUserID, collectionID, eventType string, eventData map[string]any) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO maps.audit_events (actor_user_id, collection_id, event_type, event_data)
		VALUES ($1::uuid, $2::uuid, $3, $4)
	`, actorUserID, collectionID, eventType, eventData)
	return err
}

func normalizeMemberRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "editor":
		return "editor"
	case "viewer":
		return "viewer"
	default:
		return ""
	}
}

func validateCollectionText(name, description string) error {
	if name == "" || len([]rune(name)) > 160 || len([]rune(description)) > 4000 {
		return ErrInvalidInput
	}
	return nil
}

func validateCollectionItem(provider, providerPlaceID, name, address, note string, latitude, longitude float64) error {
	if provider == "" || len([]rune(provider)) > 80 || len([]rune(providerPlaceID)) > 500 || name == "" || len([]rune(name)) > 240 || len([]rune(address)) > 1000 || len([]rune(note)) > 4000 || !validCoordinates(latitude, longitude) {
		return ErrInvalidInput
	}
	return nil
}

func validCoordinates(latitude, longitude float64) bool {
	return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

type collectionItemScanner interface {
	Scan(dest ...any) error
}

func scanCollectionItem(scanner collectionItemScanner, item *CollectionItem) error {
	return scanner.Scan(
		&item.ID,
		&item.CollectionID,
		&item.CreatedByUserID,
		&item.Provider,
		&item.ProviderPlaceID,
		&item.Name,
		&item.Address,
		&item.Latitude,
		&item.Longitude,
		&item.Note,
		&item.SortKey,
		&item.Revision,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
}
