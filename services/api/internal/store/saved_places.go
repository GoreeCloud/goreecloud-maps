package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SavedPlace struct {
	ID              string    `json:"id"`
	Provider        string    `json:"provider"`
	ProviderPlaceID string    `json:"providerPlaceId,omitempty"`
	Name            string    `json:"name"`
	Address         string    `json:"address,omitempty"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Note            string    `json:"note,omitempty"`
	Revision        int64     `json:"revision"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateSavedPlaceInput struct {
	Provider        string
	ProviderPlaceID string
	Name            string
	Address         string
	Latitude        float64
	Longitude       float64
	Note            string
}

type UpdateSavedPlaceInput struct {
	Name             *string
	Address          *string
	Latitude         *float64
	Longitude        *float64
	Note             *string
	ExpectedRevision int64
}

func (s *Store) ListSavedPlaces(ctx context.Context, userID string) ([]SavedPlace, error) {
	places := make([]SavedPlace, 0)
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT
				id::text,
				provider,
				COALESCE(provider_place_id, ''),
				name,
				COALESCE(address, ''),
				ST_Y(position::geometry),
				ST_X(position::geometry),
				COALESCE(note, ''),
				revision,
				created_at,
				updated_at
			FROM maps.saved_places
			ORDER BY updated_at DESC, id
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var place SavedPlace
			if err := rows.Scan(
				&place.ID,
				&place.Provider,
				&place.ProviderPlaceID,
				&place.Name,
				&place.Address,
				&place.Latitude,
				&place.Longitude,
				&place.Note,
				&place.Revision,
				&place.CreatedAt,
				&place.UpdatedAt,
			); err != nil {
				return err
			}
			places = append(places, place)
		}
		return rows.Err()
	})
	return places, err
}

func (s *Store) CreateSavedPlace(ctx context.Context, userID string, input CreateSavedPlaceInput) (SavedPlace, error) {
	input.Provider = strings.TrimSpace(input.Provider)
	input.ProviderPlaceID = strings.TrimSpace(input.ProviderPlaceID)
	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)
	input.Note = strings.TrimSpace(input.Note)
	if !validSavedPlaceText(input.Provider, input.ProviderPlaceID, input.Name, input.Address, input.Note) || !validCoordinates(input.Latitude, input.Longitude) {
		return SavedPlace{}, ErrInvalidInput
	}

	var place SavedPlace
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			INSERT INTO maps.saved_places (
				owner_user_id, provider, provider_place_id, name, address, position, note
			)
			VALUES (
				$1::uuid, $2, NULLIF($3, ''), $4, NULLIF($5, ''),
				ST_SetSRID(ST_MakePoint($7, $6), 4326)::geography,
				NULLIF($8, '')
			)
			RETURNING
				id::text, provider, COALESCE(provider_place_id, ''), name, COALESCE(address, ''),
				ST_Y(position::geometry), ST_X(position::geometry), COALESCE(note, ''),
				revision, created_at, updated_at
		`, userID, input.Provider, input.ProviderPlaceID, input.Name, input.Address, input.Latitude, input.Longitude, input.Note).Scan(
			&place.ID,
			&place.Provider,
			&place.ProviderPlaceID,
			&place.Name,
			&place.Address,
			&place.Latitude,
			&place.Longitude,
			&place.Note,
			&place.Revision,
			&place.CreatedAt,
			&place.UpdatedAt,
		)
		return normalizeSavedPlaceDBError(err)
	})
	return place, err
}

func (s *Store) UpdateSavedPlace(ctx context.Context, userID, savedPlaceID string, input UpdateSavedPlaceInput) (SavedPlace, error) {
	if input.ExpectedRevision < 1 || (input.Name == nil && input.Address == nil && input.Latitude == nil && input.Longitude == nil && input.Note == nil) {
		return SavedPlace{}, ErrInvalidInput
	}

	var name any
	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		if len([]rune(trimmed)) < 1 || len([]rune(trimmed)) > 240 {
			return SavedPlace{}, ErrInvalidInput
		}
		name = trimmed
	}
	var address any
	if input.Address != nil {
		trimmed := strings.TrimSpace(*input.Address)
		if len([]rune(trimmed)) > 1000 {
			return SavedPlace{}, ErrInvalidInput
		}
		address = trimmed
	}
	var note any
	if input.Note != nil {
		trimmed := strings.TrimSpace(*input.Note)
		if len([]rune(trimmed)) > 4000 {
			return SavedPlace{}, ErrInvalidInput
		}
		note = trimmed
	}
	if input.Latitude != nil && (*input.Latitude < -90 || *input.Latitude > 90) {
		return SavedPlace{}, ErrInvalidInput
	}
	if input.Longitude != nil && (*input.Longitude < -180 || *input.Longitude > 180) {
		return SavedPlace{}, ErrInvalidInput
	}

	var place SavedPlace
	err := s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, `
			UPDATE maps.saved_places
			SET
				name = COALESCE($2::text, name),
				address = CASE WHEN $3::text IS NULL THEN address ELSE NULLIF($3::text, '') END,
				position = CASE
					WHEN $4::double precision IS NULL AND $5::double precision IS NULL THEN position
					ELSE ST_SetSRID(
						ST_MakePoint(
							COALESCE($5::double precision, ST_X(position::geometry)),
							COALESCE($4::double precision, ST_Y(position::geometry))
						), 4326
					)::geography
				END,
				note = CASE WHEN $6::text IS NULL THEN note ELSE NULLIF($6::text, '') END,
				revision = revision + 1,
				updated_at = now()
			WHERE id = $1::uuid AND revision = $7
			RETURNING
				id::text, provider, COALESCE(provider_place_id, ''), name, COALESCE(address, ''),
				ST_Y(position::geometry), ST_X(position::geometry), COALESCE(note, ''),
				revision, created_at, updated_at
		`, savedPlaceID, name, address, input.Latitude, input.Longitude, note, input.ExpectedRevision).Scan(
			&place.ID,
			&place.Provider,
			&place.ProviderPlaceID,
			&place.Name,
			&place.Address,
			&place.Latitude,
			&place.Longitude,
			&place.Note,
			&place.Revision,
			&place.CreatedAt,
			&place.UpdatedAt,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			var exists bool
			if checkErr := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM maps.saved_places WHERE id = $1::uuid)`, savedPlaceID).Scan(&exists); checkErr != nil {
				return checkErr
			}
			if exists {
				return ErrConflict
			}
			return ErrNotFound
		}
		return normalizeSavedPlaceDBError(err)
	})
	return place, err
}

func (s *Store) DeleteSavedPlace(ctx context.Context, userID, savedPlaceID string, expectedRevision int64) error {
	if expectedRevision < 1 {
		return ErrInvalidInput
	}
	return s.withUserTx(ctx, userID, func(tx pgx.Tx) error {
		command, err := tx.Exec(ctx, `DELETE FROM maps.saved_places WHERE id = $1::uuid AND revision = $2`, savedPlaceID, expectedRevision)
		if err != nil {
			return err
		}
		if command.RowsAffected() == 1 {
			return nil
		}
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM maps.saved_places WHERE id = $1::uuid)`, savedPlaceID).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return ErrConflict
		}
		return ErrNotFound
	})
}

func validSavedPlaceText(provider, providerPlaceID, name, address, note string) bool {
	return len([]rune(provider)) >= 1 && len([]rune(provider)) <= 80 &&
		len([]rune(providerPlaceID)) <= 512 &&
		len([]rune(name)) >= 1 && len([]rune(name)) <= 240 &&
		len([]rune(address)) <= 1000 &&
		len([]rune(note)) <= 4000
}

func validCoordinates(latitude, longitude float64) bool {
	return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

func normalizeSavedPlaceDBError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return err
}
