package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/dead-letter/dead-letter-data/internal/pb"
	"github.com/dead-letter/dead-letter-data/migrations"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	// Read environment
	dev := os.Getenv("APP_ENV") == "development"
	dsn := os.Getenv("DATABASE_URL")

	if dev {
		fmt.Println("DEVELOPMENT MODE")
	}

	// Run migartions
	migrations.Run(dsn, dev)

	// Open database pool
	pool, err := openPool(dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Create server
	srv := &server{
		models: data.New(pool),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, srv)
	grpcServer.Serve(lis)
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
