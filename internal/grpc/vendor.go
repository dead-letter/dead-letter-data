package grpc

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/gofrs/uuid/v5"
)

type VendorServiceServer struct {
	pb.UnimplementedVendorServiceServer
	models *data.Models
}

func NewVendorServiceServer(models *data.Models) *VendorServiceServer {
	return &VendorServiceServer{
		models: models,
	}
}

func (s *VendorServiceServer) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (*pb.VendorResponse, error) {
	r, err := s.models.Vendor.Create(uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}

func (s *VendorServiceServer) ReadVendorRequest(ctx context.Context, req *pb.ReadVendorRequest) (*pb.VendorResponse, error) {
	r, err := s.models.Vendor.Read(uuid.FromStringOrNil(req.Id))
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}

func (s *VendorServiceServer) UpdateVendor(ctx context.Context, req *pb.UpdateVendorRequest) (*pb.VendorResponse, error) {
	r, err := data.VendorFromProto(req)
	if err != nil {
		return nil, err
	}

	err = s.models.Vendor.Update(r)
	if err != nil {
		return nil, err
	}

	return r.Proto(), nil
}
