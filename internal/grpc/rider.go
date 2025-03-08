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
	RiderService data.RiderService
}

func (s *RiderServiceServer) CreateRider(ctx context.Context, req *pb.CreateRiderRequest) (*pb.RiderResponse, error) {
	r, err := s.RiderService.Create(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}

func (s *RiderServiceServer) ReadRiderRequest(ctx context.Context, req *pb.ReadRiderRequest) (*pb.RiderResponse, error) {
	r, err := s.RiderService.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}

func (s *RiderServiceServer) UpdateRider(ctx context.Context, req *pb.UpdateRiderRequest) (*pb.RiderResponse, error) {
	r, err := pbconv.RiderFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.RiderService.Update(ctx, r)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromRider(r), nil
}
