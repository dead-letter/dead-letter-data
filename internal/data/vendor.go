package data

import (
	"context"

	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type VendorService interface {
	Create(ctx context.Context, id uuid.UUID) (*Vendor, error)
	Read(ctx context.Context, id uuid.UUID) (*Vendor, error)
	Update(ctx context.Context, v *Vendor) error
}

type Vendor struct {
	ID      uuid.UUID
	Version int32
}

func (v *Vendor) Proto() *pb.VendorResponse {
	return &pb.VendorResponse{
		Id:      v.ID.String(),
		Version: v.Version,
	}
}

func VendorFromProto(req *pb.UpdateVendorRequest) (*Vendor, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	v := &Vendor{
		ID:      id,
		Version: req.Version,
	}

	return v, nil
}
