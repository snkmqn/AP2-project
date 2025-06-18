package clients

import (
	orderpb "proto/generated/ecommerce/order"
	"google.golang.org/grpc"
	"log"
)

func NewOrderClient(conn *grpc.ClientConn) orderpb.OrderServiceClient {
	client := orderpb.NewOrderServiceClient(conn)
	log.Println("OrderService gRPC client initialized")
	return client
}
