package main

import (
	pb "UserManager/usermgmt"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"math/rand"
	"net"
	"os"
)

const (
	port     = ":50051"
	jsonName = "users.json"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
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

func (usm *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParamsRequest) (*pb.UserListResponse, error) {
	log.Printf("Get Users")
	jsonBytes, err := os.ReadFile(jsonName)
	if err != nil {
		log.Fatalf("Failed read from file: %v", err.Error())
	}
	var userList = &pb.UserListResponse{}
	if err := protojson.Unmarshal(jsonBytes, userList); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err.Error())
	}
	return userList, nil
}

func (usm *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	readByres, err := os.ReadFile(jsonName)
	var (
		userList    = &pb.UserListResponse{}
		userId      = int32(rand.Intn(1000))
		createdUser = &pb.UserResponse{Name: in.GetName(), Age: in.GetAge(), Id: userId}
	)
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			userList.Users = append(userList.Users, createdUser)
			if err := writeJSONUserList(userList, jsonName); err != nil {
				log.Fatalf(err.Error())
			}
			return createdUser, nil
		} else {
			log.Fatalln("Error reading file: ", err.Error())
		}
	}
	if err := protojson.Unmarshal(readByres, userList); err != nil {
		log.Fatalf("Failed to parse user list: %v", err.Error())
	}
	userList.Users = append(userList.Users, createdUser)
	if err := writeJSONUserList(userList, jsonName); err != nil {
		log.Fatalf(err.Error())
	}

	return createdUser, nil
}

func writeJSONUserList(userList *pb.UserListResponse, fileName string) error {
	jsonBytes, err := protojson.Marshal(userList)
	if err != nil {
		return errors.New(fmt.Sprintf("JSON Marchaling failde: %v", err.Error()))
	}
	if err := os.WriteFile(fileName, jsonBytes, 0664); err != nil {
		return errors.New(fmt.Sprintf("Failed write to file: %v", err.Error()))
	}
	return nil
}

func main() {
	userMGMTServ := NewUserManagementServer()
	if err := userMGMTServ.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err.Error())
	}
}
