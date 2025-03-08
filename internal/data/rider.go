package data

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/uuid"
)

type RiderService interface {
	Create(ctx context.Context, id uuid.UUID) (*Rider, error)
	Read(ctx context.Context, id uuid.UUID) (*Rider, error)
	Update(ctx context.Context, r *Rider) error
}

type Rider struct {
	ID      uuid.UUID
	Version int32
}
