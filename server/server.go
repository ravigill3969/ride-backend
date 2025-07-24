package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-grpc-auth/pb/go-grpc-auth/pb"
	database "go-grpc-auth/server/db"
	"go-grpc-auth/utils"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)


type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	_, err = database.DB.Exec(`INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`, req.Username, req.Email, hashedPassword)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &pb.AuthResponse{Message: "registered successfully"}, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("env error %s", err.Error())
		return
	}

	_, err = database.ConnectDB()
	if err != nil {
		log.Printf("Database error %s", err.Error())
		return
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// 3. Instantiate your new server struct and register it.
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	log.Println("gRPC server is running at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
