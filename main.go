package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dead-letter/dead-letter-data/pb"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type server struct {
	pb.UnimplementedDataServiceServer
	db *pgxpool.Pool
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse

	user.Id = 1
	user.Email = "johndoe@gmail.com"
	user.Name = "John Doe"

	return &user, nil
}

func newServer(db *pgxpool.Pool) *server {
	return &server{db: db}
}

func main() {
	dsn := os.Getenv("DATABASE_URL")

	fmt.Println("DSN:", dsn)

	// Run migrations
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	db.Close()

	// Open database pool
	pool, err := openPool(dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, newServer(pool))
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
