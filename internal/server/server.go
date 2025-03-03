package server

import (
	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pb.UnimplementedDataServiceServer
	models data.Models
}

// Create new data service server
func New(pool *pgxpool.Pool) *Server {
	return &Server{
		models: data.New(pool),
	}
}
