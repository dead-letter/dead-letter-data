package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

func ProtoFromRider(u *data.Rider) *pb.Rider {
	return &pb.Rider{
		Id:      u.ID.String(),
		Version: u.Version,
	}
}
