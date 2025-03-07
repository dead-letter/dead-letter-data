package data

import (
	"github.com/gofrs/uuid/v5"
)

type RiderService interface {
	Create(ID uuid.UUID) (*Rider, error)
	Read(ID uuid.UUID) (*Rider, error)
	Update(Rider *Rider) error
}

type Rider struct {
	ID      uuid.UUID
	Version int
}
