package main

import (
	pb "UserManager/usermgmt"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("didn connect: %v", err.Error())
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var newUser = make(map[string]int32)
	newUser["Alice"] = 43
	newUser["Bob"] = 30

	//Create New User
	for name, age := range newUser {
		r, err := c.CreateNewUser(ctx, &pb.NewUserRequest{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not crate user %v", err.Error())
		}
		log.Printf(`User Details:
	NAME: %s
	AGE: %d
	ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}

	//Get List Users
	params := &pb.GetUsersParamsRequest{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not rretrieve users: %v", err.Error())
	}
	log.Printf("\nUser LIST:\n")
	fmt.Printf("r.GetUsers(): %v\n", r.Users)
}
