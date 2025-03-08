package data

import (
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type VendorModel interface {
	Create(ID uuid.UUID) (*Vendor, error)
	Read(ID uuid.UUID) (*Vendor, error)
	Update(Vendor *Vendor) error
}

type Vendor struct {
	ID      uuid.UUID
	Version int32
}

func (r *Vendor) Proto() *pb.VendorResponse {
	return &pb.VendorResponse{
		Id:      r.ID.String(),
		Version: r.Version,
	}
}

func VendorFromProto(req *pb.UpdateVendorRequest) (*Vendor, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	r := &Vendor{
		ID:      id,
		Version: req.Version,
	}

	return r, nil
}
