package data

import (
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type RiderModel interface {
	Create(ID uuid.UUID) (*Rider, error)
	Read(ID uuid.UUID) (*Rider, error)
	Update(Rider *Rider) error
}

type Rider struct {
	ID      uuid.UUID
	Version int32
}

func (r *Rider) Proto() *pb.RiderResponse {
	return &pb.RiderResponse{
		Id:      r.ID.String(),
		Version: r.Version,
	}
}

func RiderFromProto(req *pb.UpdateRiderRequest) (*Rider, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	r := &Rider{
		ID:      id,
		Version: req.Version,
	}

	return r, nil
}
