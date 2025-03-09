package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoFromUser(u *data.User) *pb.User {
	return &pb.User{
		Id:        u.ID.String(),
		Version:   u.Version,
		CreatedAt: timestamppb.New(u.CreatedAt),
		Email:     u.Email,
	}
}
