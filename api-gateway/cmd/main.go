package main

import (
	"api-gateway/handlers"
	"api-gateway/router"
	"api-gateway/internal/delivery/grpc/middleware"
	middleware2 "api-gateway/internal/delivery/http/middleware"
	"proto/generated/ecommerce/inventory"
	"proto/generated/ecommerce/order"
	"proto/generated/ecommerce/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	connUser, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(middleware.ClientInterceptor("API Gateway")))
	if err != nil {
		log.Fatalf("Could not connect to User Service: %v", err)
	}
	defer connUser.Close()

	connProduct, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(middleware.ClientInterceptor("API Gateway")))
	if err != nil {
		log.Fatalf("Could not connect to Product Service: %v", err)
	}

	defer connProduct.Close()

	connOrder, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(middleware.ClientInterceptor("API Gateway")))
	if err != nil {
		log.Fatalf("Could not connect to Order Service: %v", err)
	}
	defer connOrder.Close()

	productClient := inventorypb.NewInventoryServiceClient(connProduct)
	orderClient := orderpb.NewOrderServiceClient(connOrder)
	userClient := userpb.NewUserServiceClient(connUser)

	productHandler := handlers.NewProductHandler(productClient)
	orderHandler := handlers.NewOrderHandler(orderClient)
	userHandler := handlers.NewUserHandler(userClient)

	r := gin.New()
	r.Use(middleware2.LoggerMiddleware("API Gateway"))
	r.Use(gin.Recovery())

	router.SetupRouter(r, productHandler, orderHandler, userHandler)

	log.Println("API Gateway is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
