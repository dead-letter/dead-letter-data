package pg

import (
	"context"
	"errors"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RiderService struct {
	Pool *pgxpool.Pool
}

func (s *RiderService) Create(ctx context.Context, id uuid.UUID) (*data.Rider, error) {
	r := data.Rider{
		ID: id,
	}

	sql := `
		INSERT INTO rider_ (id_)
		VALUES($1)
		RETURNING version_;`

	err := s.Pool.QueryRow(ctx, sql, r.ID).Scan(
		&r.Version,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *RiderService) Read(ctx context.Context, id uuid.UUID) (*data.Rider, error) {
	var r data.Rider

	sql := `
		SELECT id_, version_
		FROM rider_ WHERE id_ = $1;`

	err := s.Pool.QueryRow(ctx, sql, id).Scan(
		&r.ID,
		&r.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &r, nil
}

func (s *RiderService) Update(ctx context.Context, r *data.Rider) error {
	sql := `
		UPDATE rider_ 
        SET version_ = version_ + 1
        WHERE id_ = $1 AND version_ = $2
        RETURNING version_;`

	args := []any{
		r.ID,
		r.Version,
	}

	err := s.Pool.QueryRow(ctx, sql, args...).Scan(
		&r.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return data.ErrEditConflict
		case pgErrCode(err) == pgerrcode.UniqueViolation:
			return data.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
