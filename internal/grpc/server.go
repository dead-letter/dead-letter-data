package grpc

import (
	"net"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	Addr string
	DB   *data.DB
}

func NewServer(addr string, db *data.DB) *Server {
	return &Server{
		Addr: addr,
		DB:   db,
	}
}

func (s *Server) ListenAndServe() error {
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &UserServiceServer{DB: s.DB})
	pb.RegisterRiderServiceServer(grpcServer, &RiderServiceServer{DB: s.DB})
	pb.RegisterVendorServiceServer(grpcServer, &VendorServiceServer{DB: s.DB})

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
