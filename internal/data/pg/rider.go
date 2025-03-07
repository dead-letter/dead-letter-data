package pg

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RiderModel struct {
	pool *pgxpool.Pool
}

type Rider struct {
	UserID uuid.UUID
}

func (r *Rider) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required))
}

func (m RiderModel) New(userID uuid.UUID) (*Rider, error) {
	r := &Rider{UserID: userID}

	err := m.Insert(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m RiderModel) Insert(r *Rider) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO rider_ (id_)
		VALUES($1);`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err = m.pool.Exec(ctx, sql, r.UserID)
	if err != nil {
		return err
	}

	return nil
}
