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
	UserService data.UserService
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := s.UserService.Create(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) ReadUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.UserResponse, error) {
	u, err := s.UserService.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) ReadUserWithEmailRequest(ctx context.Context, req *pb.ReadUserWithEmailRequest) (*pb.UserResponse, error) {
	u, err := s.UserService.ReadWithEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) ReadUserWithCredentialsRequest(ctx context.Context, req *pb.ReadUserWithCredentialsRequest) (*pb.UserResponse, error) {
	u, err := s.UserService.ReadWithCredentials(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromUser(u), nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := pbconv.UserFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.UserService.Update(ctx, u)
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

	err = s.UserService.Delete(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
