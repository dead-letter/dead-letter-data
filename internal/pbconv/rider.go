package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

func ProtoFromRider(u *data.Rider) *pb.RiderResponse {
	return &pb.RiderResponse{
		Id:      u.ID.String(),
		Version: u.Version,
	}
}

func RiderFromProto(req *pb.UpdateRiderRequest) (*data.Rider, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	u := &data.Rider{
		ID:      id,
		Version: req.Version,
	}

	return u, nil
}
