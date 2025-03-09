package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/pbconv"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	DB *data.DB
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	u, err := s.DB.Users.Create(ctx, req.Email, req.PasswordHash)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) ReadUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.User, error) {
	u, err := s.DB.Users.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) CheckUserExists(ctx context.Context, req *pb.CheckUserExistsRequest) (*pb.CheckUserExistsResponse, error) {
	exists, err := s.DB.Users.ExistsWithEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserExistsResponse{Exists: exists}, nil
}

func (s *UserServiceServer) ReadUserWithCredentials(ctx context.Context, req *pb.ReadUserWithCredentialsRequest) (*pb.User, error) {
	u, err := s.DB.Users.ReadWithCredentials(ctx, req.Email, req.PasswordHash)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	u := &data.User{
		ID:           id,
		Version:      req.Version,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
	}

	err = s.DB.Users.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.DB.Users.Delete(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
