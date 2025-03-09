package pbconv

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

func ProtoFromVendor(v *data.Vendor) *pb.Vendor {
	return &pb.Vendor{
		Id:      v.ID.String(),
		Version: v.Version,
	}
}
