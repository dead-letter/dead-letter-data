package rpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type RiderServiceServer struct {
	pb.UnimplementedRiderServiceServer
	models *data.Models
}

func NewRiderServiceServer(models *data.Models) *RiderServiceServer {
	return &RiderServiceServer{
		models: models,
	}
}

func (s *RiderServiceServer) CreateRider(ctx context.Context, req *pb.CreateRiderRequest) (*pb.RiderResponse, error) {
	r, err := s.models.Rider.Create(uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}

func (s *RiderServiceServer) ReadRiderRequest(ctx context.Context, req *pb.ReadRiderRequest) (*pb.RiderResponse, error) {
	r, err := s.models.Rider.Read(uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}

func (s *RiderServiceServer) UpdateRider(ctx context.Context, req *pb.UpdateRiderRequest) (*pb.RiderResponse, error) {
	r, err := data.RiderFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.models.Rider.Update(r)
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}
