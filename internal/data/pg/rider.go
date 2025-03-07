package pg

import (
	"context"
	"errors"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RiderModel struct {
	pool *pgxpool.Pool
}

func NewRiderModel(pool *pgxpool.Pool) RiderModel {
	return RiderModel{
		pool: pool,
	}
}

func (s RiderModel) Create(id uuid.UUID) (*data.Rider, error) {
	r := data.Rider{
		ID: id,
	}

	sql := `
		INSERT INTO rider_ (id_)
		VALUES($1)
		RETURNING version_;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, r.ID).Scan(
		&r.Version,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s RiderModel) Read(id uuid.UUID) (*data.Rider, error) {
	var r data.Rider

	sql := `
		SELECT id_, version_
		FROM rider_ WHERE id_ = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, id).Scan(
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

func (s RiderModel) Update(r *data.Rider) error {
	sql := `
		UPDATE rider_ 
        SET version_ = version_ + 1
        WHERE id_ = $1 AND version_ = $2
        RETURNING version_;`

	args := []any{
		r.ID,
		r.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, args...).Scan(
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
