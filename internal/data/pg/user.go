package pg

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	pool *pgxpool.Pool
}

func NewUserModel(pool *pgxpool.Pool) UserModel {
	return UserModel{pool: pool}
}

func (s UserModel) Create(email, password string) (*data.User, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	u := data.User{
		Email:        email,
		PasswordHash: []byte(hash),
	}

	sql := `
		INSERT INTO user_ (email_, password_hash_)
		VALUES($1, $2)
		RETURNING id_, version_, created_at_;`

	args := []any{u.Email, u.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err = s.pool.QueryRow(ctx, sql, args...).Scan(
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

func (s UserModel) Read(id uuid.UUID) (*data.User, error) {
	var u data.User

	sql := `
		SELECT id_, version_, created_at_, email_, password_hash_
		FROM user_ WHERE id_ = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, id).Scan(
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

func (s UserModel) ReadWithEmail(email string) (*data.User, error) {
	var u data.User

	sql := `
		SELECT id_, version_, created_at_, email_, password_hash_
		FROM user_ WHERE email_ = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, email).Scan(
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

func (s UserModel) ReadWithCredentials(email, password string) (*data.User, error) {
	u, err := s.ReadWithEmail(email)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, string(u.PasswordHash))
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, data.ErrInvalidCredentials
	}

	return u, nil
}

func (s UserModel) Update(u *data.User) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err := s.pool.QueryRow(ctx, sql, args...).Scan(
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

func (s UserModel) Delete(id uuid.UUID) error {
	sql := `
		DELETE FROM user_
		WHERE id_ = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	res, err := s.pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
