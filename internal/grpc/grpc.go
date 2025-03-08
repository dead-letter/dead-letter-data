package grpc

import (
	"net"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	Addr   string
	Models *data.Models
}

func (s *Server) ListenAndServe() error {
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, NewUserServiceServer(s.Models))
	pb.RegisterRiderServiceServer(grpcServer, NewRiderServiceServer(s.Models))

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
