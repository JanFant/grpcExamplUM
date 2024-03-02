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

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		userList: &pb.UserListResponse{},
	}
}

func (serv *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, serv)

	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	userList *pb.UserListResponse
}

func (usm *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	var userId int32 = int32(rand.Intn(1000))
	createdUser := &pb.UserResponse{Name: in.GetName(), Age: in.GetAge(), Id: userId}
	usm.userList.Users = append(usm.userList.Users, createdUser)
	return createdUser, nil
}

func (usm *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParamsRequest) (*pb.UserListResponse, error) {
	log.Printf("Get Users")
	return usm.userList, nil
}

func main() {
	userMGMTServ := NewUserManagementServer()
	if err := userMGMTServ.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err.Error())
	}
}
