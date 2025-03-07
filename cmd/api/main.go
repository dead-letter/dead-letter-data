package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/data/pg"
	"github.com/dead-letter/dead-letter-data/internal/server"
	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)

func main() {
	var dev bool

	// CLI flags
	flag.BoolVar(&dev, "dev", false, "Development mode")
	flag.Parse()

	// Read environment
	dsn := os.Getenv("DATABASE_URL")

	// Logger
	h := newSlogHandler(dev)
	logger := slog.New(h)

	// Open database pool
	pool, err := pg.OpenPool(dsn)
	if err != nil {
		fatal(logger, err)
	}
	defer pool.Close()

	// Run migrations
	db := stdlib.OpenDBFromPool(pool)
	migrations.Up(db)
	if dev {
		migrations.Reset(db)
	}
	db.Close()

	// Create models
	models := &data.Models{
		User:  pg.NewUserService(pool),
		Rider: pg.NewRiderService(pool),
	}

	// Run gRPC server
	srv := server.New(models)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fatal(logger, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, srv)

	logger.Info("starting gRPC server", slog.String("addr", lis.Addr().String()))
	grpcServer.Serve(lis)
	if err != nil {
		fatal(logger, err)
	}
}

func newSlogHandler(dev bool) slog.Handler {
	if dev {
		// Development text hanlder
		return tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		})
	}

	// Production use JSON handler with default opts
	return slog.NewJSONHandler(os.Stdout, nil)
}

func fatal(logger *slog.Logger, err error) {
	logger.Error("fatal", slog.Any("err", err))
	os.Exit(1)
}
