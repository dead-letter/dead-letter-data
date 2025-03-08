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

type VendorService struct {
	Pool *pgxpool.Pool
}

func (s *VendorService) Create(ctx context.Context, id uuid.UUID) (*data.Vendor, error) {
	v := data.Vendor{
		ID: id,
	}

	sql := `
		INSERT INTO vendor_ (id_)
		VALUES($1)
		RETURNING version_;`

	err := s.Pool.QueryRow(ctx, sql, v.ID).Scan(
		&v.Version,
	)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *VendorService) Read(ctx context.Context, id uuid.UUID) (*data.Vendor, error) {
	var v data.Vendor

	sql := `
		SELECT id_, version_
		FROM vendor_ WHERE id_ = $1;`

	err := s.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID,
		&v.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &v, nil
}

func (s *VendorService) Update(ctx context.Context, v *data.Vendor) error {
	sql := `
		UPDATE vendor_ 
        SET version_ = version_ + 1
        WHERE id_ = $1 AND version_ = $2
        RETURNING version_;`

	args := []any{
		v.ID,
		v.Version,
	}

	err := s.Pool.QueryRow(ctx, sql, args...).Scan(
		&v.Version,
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
