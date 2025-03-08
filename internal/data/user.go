package data

import (
	"context"
	"time"

	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	Create(ctx context.Context, email, password string) (*User, error)
	Read(ctx context.Context, id uuid.UUID) (*User, error)
	ReadWithEmail(ctx context.Context, email string) (*User, error)
	ReadWithCredentials(ctx context.Context, email, password string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type User struct {
	ID           uuid.UUID
	Version      int32
	CreatedAt    time.Time
	Email        string
	PasswordHash []byte
}

func (u *User) Proto() *pb.UserResponse {
	return &pb.UserResponse{
		Id:        u.ID.String(),
		Version:   u.Version,
		CreatedAt: timestamppb.New(u.CreatedAt),
		Email:     u.Email,
	}
}

func UserFromProto(req *pb.UpdateUserRequest) (*User, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:           id,
		Version:      req.Version,
		Email:        req.Email,
		PasswordHash: []byte(req.PasswordHash),
	}

	return u, nil
}
