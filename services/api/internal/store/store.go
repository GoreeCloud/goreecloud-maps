package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	if input.Name == "" || len([]rune(input.Name)) > 160 {
		return Collection{}, errors.New("collection name must contain between 1 and 160 characters")
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
