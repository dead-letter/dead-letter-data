package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/uuid"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

func ProtoFromVendor(v *data.Vendor) *pb.VendorResponse {
	return &pb.VendorResponse{
		Id:      v.ID.String(),
		Version: v.Version,
	}
}

func VendorFromProto(req *pb.UpdateVendorRequest) (*data.Vendor, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	v := &data.Vendor{
		ID:      id,
		Version: req.Version,
	}

	return v, nil
}
