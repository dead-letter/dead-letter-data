package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/pbconv"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

type RiderServiceServer struct {
	pb.UnimplementedRiderServiceServer
	DB *data.DB
}

func (s *RiderServiceServer) CreateRider(ctx context.Context, req *pb.CreateRiderRequest) (*pb.Rider, error) {
	r, err := s.DB.Riders.Create(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}

func (s *RiderServiceServer) ReadRiderRequest(ctx context.Context, req *pb.ReadRiderRequest) (*pb.Rider, error) {
	r, err := s.DB.Riders.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}

func (s *RiderServiceServer) UpdateRider(ctx context.Context, req *pb.UpdateRiderRequest) (*pb.Rider, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	r := &data.Rider{
		ID:      id,
		Version: req.Version,
	}

	err = s.DB.Riders.Update(ctx, r)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}
