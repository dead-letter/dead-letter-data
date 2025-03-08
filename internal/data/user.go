package data

import (
	"context"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/uuid"
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
