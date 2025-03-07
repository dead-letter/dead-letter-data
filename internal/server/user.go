package server

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u, err := srv.users.Create(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (srv *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	u, err := srv.users.ReadWithEmail(req.Email)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (srv *Server) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.UserResponse, error) {
	u, err := srv.users.ReadWithEmailAndPassword(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (srv *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := data.UserFromProto(req)
	if err != nil {
		return nil, err
	}

	err = srv.users.Update(u)
	if err != nil {
		return nil, err
	}

	res := &pb.UserResponse{
		User: u.Proto(),
	}

	return res, nil
}

func (srv *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	err = srv.users.Delete(userID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
