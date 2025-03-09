package data

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/uuid"
)

type VendorRepository interface {
	Create(ctx context.Context, id uuid.UUID) (*Vendor, error)
	Read(ctx context.Context, id uuid.UUID) (*Vendor, error)
	Update(ctx context.Context, v *Vendor) error
}

type Vendor struct {
	ID      uuid.UUID
	Version int32
}
