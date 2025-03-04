package server

import (
	"log/slog"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pb.UnimplementedDataServiceServer
	log    *slog.Logger
	models data.Models
}

// Create new data service server
func New(logger *slog.Logger, pool *pgxpool.Pool) *Server {
	return &Server{
		log:    logger,
		models: data.New(pool),
	}
}
