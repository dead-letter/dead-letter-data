package grpc

import (
	"net"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	Addr          string
	UserService   data.UserService
	RiderService  data.RiderService
	VendorService data.VendorService
}

func (srv *Server) ListenAndServe() error {
	s := grpc.NewServer()

	userServiceServer := &UserServiceServer{UserService: srv.UserService}
	riderServiceServer := &RiderServiceServer{RiderService: srv.RiderService}
	vendorServiceServer := &VendorServiceServer{VendorService: srv.VendorService}

	pb.RegisterUserServiceServer(s, userServiceServer)
	pb.RegisterRiderServiceServer(s, riderServiceServer)
	pb.RegisterVendorServiceServer(s, vendorServiceServer)

	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	err = s.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
