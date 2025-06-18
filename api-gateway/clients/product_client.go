package clients

import (
	inventorypb "proto/generated/ecommerce/inventory"
	"google.golang.org/grpc"
	"log"
)

func NewProductClient(conn *grpc.ClientConn) inventorypb.InventoryServiceClient {
	client := inventorypb.NewInventoryServiceClient(conn)
	log.Println("ProductService gRPC client initialized")
	return client
}
