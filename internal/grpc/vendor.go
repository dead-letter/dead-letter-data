package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/pbconv"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

type VendorServiceServer struct {
	pb.UnimplementedVendorServiceServer
	DB *data.DB
}

func (s *VendorServiceServer) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (*pb.Vendor, error) {
	v, err := s.DB.Vendors.Create(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromVendor(v), nil
}

func (s *VendorServiceServer) ReadVendorRequest(ctx context.Context, req *pb.ReadVendorRequest) (*pb.Vendor, error) {
	v, err := s.DB.Vendors.Read(ctx, uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromVendor(v), nil
}

func (s *VendorServiceServer) UpdateVendor(ctx context.Context, req *pb.UpdateVendorRequest) (*pb.Vendor, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	v := &data.Vendor{
		ID:      id,
		Version: req.Version,
	}

	err = s.DB.Vendors.Update(ctx, v)
	if err != nil {
		return nil, err
	}

	return pbconv.ProtoFromVendor(v), nil
}
