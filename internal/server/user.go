package server

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.Create(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.ReadWithEmail(req.Email)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (s *Server) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.UserResponse, error) {
	u, err := s.models.User.ReadWithCredentials(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := data.UserFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.models.User.Update(u)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
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
