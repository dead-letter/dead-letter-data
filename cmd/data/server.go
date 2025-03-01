package main

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/pb"
)

type server struct {
	pb.UnimplementedDataServiceServer
	models data.Models
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	u, err := s.models.User.GetWithEmail(req.Email)
	if err != nil {
		return nil, err
	}

	var res pb.GetUserResponse

	res.Id = u.ID.String()
	res.CreatedAt = u.CreatedAt.String()

	return &res, nil
}
