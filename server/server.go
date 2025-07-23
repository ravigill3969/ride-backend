package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-grpc-auth/pb/go-grpc-auth/pb"

	"google.golang.org/grpc"
)

type authServer struct {
	pb.UnimplementedAuthServiceServer
	users map[string]*User
}

type User struct {
	Username     string
	Email        string
	Password     string
	SessionToken string
}

func newServer() *authServer {
	return &authServer{
		users: make(map[string]*User),
	}
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if _, exists := s.users[req.Username]; exists {
		return nil, fmt.Errorf("user already exists")
	}
	s.users[req.Username] = &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	return &pb.AuthResponse{Message: "registered successfully"}, nil
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, exists := s.users[req.Username]
	if !exists || user.Password != req.Password {
		return nil, fmt.Errorf("invalid credentials")
	}
	user.SessionToken = "session_" + req.Username
	return &pb.AuthResponse{Message: "login successful", SessionToken: user.SessionToken}, nil
}

func (s *authServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.GenericResponse, error) {
	user, exists := s.users[req.Username]
	if !exists || user.Password != req.OldPassword {
		return nil, fmt.Errorf("invalid credentials")
	}
	user.Password = req.NewPassword
	return &pb.GenericResponse{Message: "password changed"}, nil
}

func (s *authServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.GenericResponse, error) {
	user, exists := s.users[req.Username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	user.SessionToken = ""
	return &pb.GenericResponse{Message: "logged out"}, nil
}

func (s *authServer) SendCookie(ctx context.Context, req *pb.SendCookieRequest) (*pb.GenericResponse, error) {
	user, exists := s.users[req.Username]
	if !exists || user.SessionToken != req.SessionToken {
		return nil, fmt.Errorf("invalid session")
	}
	return &pb.GenericResponse{Message: "cookie accepted"}, nil
}

func (s *authServer) GetUserData(ctx context.Context, req *pb.UserDataRequest) (*pb.UserDataResponse, error) {
	user, exists := s.users[req.Username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return &pb.UserDataResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, newServer())
	log.Println("gRPC server is running at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
