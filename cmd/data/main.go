package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/server"
	"github.com/dead-letter/dead-letter-data/migrations"
	"github.com/dead-letter/dead-letter-data/pkg/pb"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)

func main() {
	// Read environment
	dev := os.Getenv("APP_ENV") == "development"
	dsn := os.Getenv("DATABASE_URL")

	// Logger
	h := newSlogHandler(dev)
	logger := slog.New(h)

	// Run migartions
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fatal(logger, err)
	}

	migrations.Up(db)
	if dev {
		migrations.Reset(db)
	}

	db.Close()

	// Open database pool
	pool, err := openPool(dsn)
	if err != nil {
		fatal(logger, err)
	}
	defer pool.Close()

	// Run gRPC server
	srv := server.New(logger, pool)

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

func openPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, err
}
