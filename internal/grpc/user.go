package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	models *data.Models
}

func NewUserServiceServer(models *data.Models) *UserServiceServer {
	return &UserServiceServer{
		models: models,
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.Create(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return u.Proto(), nil
}

func (s *UserServiceServer) ReadUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.Read(uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return u.Proto(), nil
}

func (s *UserServiceServer) ReadUserWithEmailRequest(ctx context.Context, req *pb.ReadUserWithEmailRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.ReadWithEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return u.Proto(), nil
}

func (s *UserServiceServer) ReadUserWithCredentialsRequest(ctx context.Context, req *pb.ReadUserWithCredentialsRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.ReadWithCredentials(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return u.Proto(), nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := data.UserFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.models.User.Update(u)
	if err != nil {
		return nil, err
	}

	return u.Proto(), nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.models.User.Delete(userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
