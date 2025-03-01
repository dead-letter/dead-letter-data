package main

import (
	"context"
	"log"
	"net"

	"github.com/dead-letter/dead-letter-data/pb"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDataServiceServer
	db *pgx.Conn
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse

	user.Id = 1
	user.Email = "johndoe@gmail.com"
	user.Name = "John Doe"

	return &user, nil
}

func newServer(db *pgx.Conn) *server {
	return &server{db: db}
}

func main() {
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Unable to connect to database: %v\n", err)
	// }
	// defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataServiceServer(grpcServer, newServer(nil))
	grpcServer.Serve(lis)
}
