package data

import (
	"github.com/gofrs/uuid/v5"
)

type RiderService interface {
	Create(userID uuid.UUID) (*Rider, error)
	Read(userID uuid.UUID) (*Rider, error)
	Update(Rider *Rider) error
	Delete(userID uuid.UUID) error
}

type Rider struct {
	User *User
}
