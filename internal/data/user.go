package data

import (
	"time"

	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserModel interface {
	Create(email, password string) (*User, error)
	Read(id uuid.UUID) (*User, error)
	ReadWithEmail(email string) (*User, error)
	ReadWithCredentials(email, password string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
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
