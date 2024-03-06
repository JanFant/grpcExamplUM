package main

import (
	pb "UserManager/usermgmt"
	"UserManager/usermgmtServer/db"
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"net"
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
	conn *pgx.Conn
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
	var userList = &pb.UserListResponse{}
	rows, err := usm.conn.Query(context.Background(), `SELECT * FROM "testDB".public.users`)
	if err != nil {
		log.Printf("cannt take data %v", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &pb.UserResponse{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, user)
	}
	return userList, nil
}

func (usm *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.UserResponse, error) {
	log.Printf("Received: %v", in.GetName())
	var (
		createdUser = &pb.UserResponse{Name: in.GetName(), Age: in.GetAge()}
	)
	tx, err := usm.conn.Begin(context.Background())
	if err != nil {
		log.Printf("conn.Begin failed: %v", err.Error())
		return nil, err
	}
	err = tx.QueryRow(context.Background(), `INSERT INTO "testDB".public.users (name, age) values ($1, $2) returning id`,
		createdUser.Name, createdUser.Age).Scan(&createdUser.Id)
	if err != nil {
		log.Printf("Insert error: %v", err.Error())
		return nil, err
	}
	_ = tx.Commit(context.Background())

	return createdUser, nil
}

func main() {
	userMGMTServ := NewUserManagementServer()

	conn, err := db.ConnectPDB()
	if err != nil {
		log.Fatalf("Cound not connect DB")
	}
	defer conn.Close(context.Background())
	userMGMTServ.conn = conn
	if err := userMGMTServ.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err.Error())
	}
}
