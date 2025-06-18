package clients

import (
	userpb "proto/generated/ecommerce/user"
	"google.golang.org/grpc"
	"log"
)

func NewUserClient(conn *grpc.ClientConn) userpb.UserServiceClient {
	client := userpb.NewUserServiceClient(conn)
	log.Println("UserService gRPC client initialized")
	return client
}
