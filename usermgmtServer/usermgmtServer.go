package main

import (
	pb "UserManager/usermgmt"
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	var userId int32 = int32(rand.Intn(1000))
	return &pb.UserResponse{Name: in.GetName(), Age: in.GetAge(), Id: userId}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err.Error())
	}
}
