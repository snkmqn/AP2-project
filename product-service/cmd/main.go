package main

import (
	"context"
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"product-service/internal/config"
	"product-service/internal/delivery/grpc/middleware"
	"product-service/internal/infrastructure/cache"
	"product-service/internal/infrastructure/database"
	"product-service/internal/infrastructure/logger"
	"product-service/internal/infrastructure/repositories"
	"product-service/internal/infrastructure/utils/jwt"
	grpc2 "product-service/internal/interfaces/grpc"
	repositories2 "product-service/internal/interfaces/repositories"
	"product-service/internal/migrations/migrations"
	mongo2 "product-service/internal/migrations/mongo"
	"proto/generated/ecommerce/inventory"
)

func initRepositories() (repositories2.ProductRepository, *mongo.Database, error) {
	db, err := database.ConnectMongoDB("inventory")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	mongo2.RegisterMigrations()

	if err := migrations.Run(context.Background(), db); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	productRepo := repositories.NewProductRepositoryMongo(db)
	return productRepo, db, nil
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics available on :8082/metrics")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

func main() {
	productRepo, _, err := initRepositories()
	if err != nil {
		log.Fatal(err)
	}
	secretKey := config.GetEnv("JWT_SECRET_KEY", "")

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

	stdLogger := &logger.StdLogger{}
	inventoryServer := grpc2.NewInventoryGrpcServer(productRepo, stdLogger, redisClient)
	inventorypb.RegisterInventoryServiceServer(grpcServer, inventoryServer)

	go startMetricsServer()

	grpc_prometheus.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}
	log.Println("Product Service is running on port :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
