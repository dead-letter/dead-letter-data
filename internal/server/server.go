package server

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
)

type Server struct {
	pb.UnimplementedDataServiceServer
	models *data.Models
}

func New(models *data.Models) *Server {
	return &Server{models: models}
}
