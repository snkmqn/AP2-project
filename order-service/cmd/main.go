package main

import (
	"order-service/internal/config"
	"order-service/internal/delivery/grpc/middleware"
	"order-service/internal/infrastructure/cache"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/logger"
	"order-service/internal/infrastructure/repositories"
	"order-service/internal/infrastructure/utils/jwt"
	"order-service/internal/infrastructure/utils/uuid"
	grpc2 "order-service/internal/interfaces/grpc"
	repositories2 "order-service/internal/interfaces/repositories"
	"order-service/internal/usecases/services"
	"proto/generated/ecommerce/order"
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func initRepositories() (repositories2.OrderRepository, repositories2.ProductRepository, error) {
	dbOrders, err := database.ConnectMongoDB("orders")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to orders DB: %v", err)
	}

	dbInventory, err := database.ConnectMongoDB("inventory")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to inventory DB: %v", err)
	}

	orderRepo := repositories.NewOrderRepositoryMongo(dbOrders)
	productRepo := repositories.NewProductRepositoryMongo(dbInventory)

	return orderRepo, productRepo, nil
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics available on :8083/metrics")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	natsURL := os.Getenv("NATS_URL")

	orderRepo, productRepo, err := initRepositories()
	if err != nil {
		log.Fatal(err)
	}
	secretKey := config.GetEnv("JWT_SECRET_KEY", "")
	stdLogger := &logger.StdLogger{}

	redisAddr := config.GetEnv("REDIS_ADDR", "")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")
	redisDB := 0
	redisClient := cache.NewRedisCache(redisAddr, redisPassword, redisDB)

	var jwtService jwt.JWTService = jwt.NewJWTService(secretKey, redisClient)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			middleware.JWTInterceptor(jwtService),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
		),
	)

	natsClient, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	uuidGen := uuid.NewUUIDService()
	priceCalculator := services.NewPriceCalculator()
	productService := services.NewProductService(productRepo, stdLogger, redisClient)

	orderService := services.NewOrderService(orderRepo, priceCalculator, uuidGen, productService, redisClient, stdLogger)
	orderServer := grpc2.NewOrderGrpcServer(stdLogger, orderService, natsClient, redisClient)
	orderpb.RegisterOrderServiceServer(grpcServer, orderServer)

	go startMetricsServer()

	grpc_prometheus.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053: %v", err)
	}
	log.Println("Order Service is running on port :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
