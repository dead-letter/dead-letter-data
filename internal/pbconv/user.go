package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoFromUser(u *data.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        u.ID.String(),
		Version:   u.Version,
		CreatedAt: timestamppb.New(u.CreatedAt),
		Email:     u.Email,
	}
}

func UserFromProto(req *pb.UpdateUserRequest) (*data.User, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	u := &data.User{
		ID:           id,
		Version:      req.Version,
		Email:        req.Email,
		PasswordHash: []byte(req.PasswordHash),
	}

	return u, nil
}
