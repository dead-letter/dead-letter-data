package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type VendorServiceServer struct {
	pb.UnimplementedVendorServiceServer
	VendorService data.VendorService
}

func (s *VendorServiceServer) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (*pb.VendorResponse, error) {
	v, err := s.VendorService.Create(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return v.Proto(), nil
}

func (s *VendorServiceServer) ReadVendorRequest(ctx context.Context, req *pb.ReadVendorRequest) (*pb.VendorResponse, error) {
	v, err := s.VendorService.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return v.Proto(), nil
}

func (s *VendorServiceServer) UpdateVendor(ctx context.Context, req *pb.UpdateVendorRequest) (*pb.VendorResponse, error) {
	v, err := data.VendorFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.VendorService.Update(ctx, v)
	if err != nil {
		return nil, err
	}

	return v.Proto(), nil
}
