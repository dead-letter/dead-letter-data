package data

import "github.com/gofrs/uuid/v5"

type VendorService interface {
	Create(userID uuid.UUID) (*Vendor, error)
	Read(userID uuid.UUID) (*Vendor, error)
	Update(Vendor *Vendor) error
	Delete(userID uuid.UUID) error
}

type Vendor struct {
	user User
}
