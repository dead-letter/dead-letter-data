package main

import (
	"context"
<<<<<<< HEAD
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/dead-letter/dead-letter-manifests/pb"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDataServiceServer
	db *pgx.Conn
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse

	sql := `
		SELECT id, name, email 
		FROM users WHERE id=$1;`

	err := s.db.QueryRow(ctx, sql, req.Id).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	return &user, nil
}

func newServer(db *pgx.Conn) *server {
	return &server{db: db}
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())


	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, newServer(conn))

	grpcServer.Serve(lis)
}
