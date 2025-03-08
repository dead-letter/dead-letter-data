package pg

import (
	"bytes"
	"context"
	"errors"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	Pool *pgxpool.Pool
}

func (s *UserService) Create(ctx context.Context, email string, passwordHash []byte) (*data.User, error) {
	u := data.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	sql := `
		INSERT INTO user_ (email_, password_hash_)
		VALUES($1, $2)
		RETURNING id_, version_, created_at_;`

	args := []any{u.Email, u.PasswordHash}

	err := s.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.ID,
		&u.Version,
		&u.CreatedAt,
	)
	if err != nil {
		switch {
		case pgErrCode(err) == pgerrcode.UniqueViolation:
			return nil, data.ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (s *UserService) Read(ctx context.Context, id uuid.UUID) (*data.User, error) {
	var u data.User

	sql := `
		SELECT id_, version_, created_at_, email_, password_hash_
		FROM user_ WHERE id_ = $1;`

	err := s.Pool.QueryRow(ctx, sql, id).Scan(
		&u.ID,
		&u.Version,
		&u.CreatedAt,
		&u.Email,
		&u.PasswordHash,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (s *UserService) ExistsWithEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM user_
			WHERE email_ = $1
		);`

	err := s.Pool.QueryRow(ctx, sql, email).Scan(&exists)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return false, nil
		default:
			return exists, err
		}
	}

	return exists, nil
}

func (s *UserService) ReadWithCredentials(ctx context.Context, email string, passwordHash []byte) (*data.User, error) {
	var u data.User

	sql := `
		SELECT id_, version_, created_at_, email_, password_hash_
		FROM user_ WHERE email_ = $1;`

	err := s.Pool.QueryRow(ctx, sql, email).Scan(
		&u.ID,
		&u.Version,
		&u.CreatedAt,
		&u.Email,
		&u.PasswordHash,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if !bytes.Equal(passwordHash, u.PasswordHash) {
		return nil, data.ErrInvalidCredentials
	}

	return &u, nil
}

func (s *UserService) Update(ctx context.Context, u *data.User) error {
	sql := `
		UPDATE user_ 
        SET email_ = $1, password_hash_ = $2, version_ = version_ + 1
        WHERE id_ = $3 AND version_ = $4
        RETURNING version_;`

	args := []any{
		u.Email,
		u.PasswordHash,
		u.ID,
		u.Version,
	}

	err := s.Pool.QueryRow(ctx, sql, args...).Scan(
		&u.Version,
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

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
		DELETE FROM user_
		WHERE id_ = $1;`

	res, err := s.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
