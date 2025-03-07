package pg

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VendorModel struct {
	pool *pgxpool.Pool
}

type Vendor struct {
	UserID uuid.UUID
}

func (r *Vendor) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required))
}

func (m VendorModel) New(userID uuid.UUID) (*Vendor, error) {
	r := &Vendor{UserID: userID}

	err := m.Insert(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m VendorModel) Insert(r *Vendor) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO vendor_ (id_)
		VALUES($1);`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err = m.pool.Exec(ctx, sql, r.UserID)
	if err != nil {
		return err
	}

	return nil
}
