package main

import (
	"consumer-service/internal/config"
	"consumer-service/internal/infrastructure/cache"
	"consumer-service/internal/infrastructure/database"
	"consumer-service/internal/infrastructure/logger"
	"consumer-service/internal/infrastructure/messaging"
	repositories2 "consumer-service/internal/infrastructure/repositories"
	"consumer-service/internal/interfaces/repositories"
	"consumer-service/internal/usecases/services"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func initRepositories() (repositories.ProductRepository, error) {
	dbInventory, err := database.ConnectMongoDB("inventory")

	if err != nil {
		return nil, fmt.Errorf("failed to connect to orders DB: %v", err)
	}

	productRepo := repositories2.NewProductRepositoryMongo(dbInventory)

	return productRepo, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	natsURL := os.Getenv("NATS_URL")
	grpcServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")

	productRepo, err := initRepositories()
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}

	stdLogger := &logger.StdLogger{}

	grpcConn, err := grpc.NewClient(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product-service: %v", err)
	}

	defer grpcConn.Close()

	redisAddr := config.GetEnv("REDIS_ADDR", "")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")
	redisDB := 0
	redisClient := cache.NewRedisCache(redisAddr, redisPassword, redisDB)

	productService := services.NewProductService(productRepo, stdLogger, redisClient)
	stockDecrementer := services.NewStockDecrementer(*productService)

	stockSubscriber, err := messaging.NewStockSubscriber(natsURL, *stockDecrementer, productRepo, grpcServerAddress)

	if err != nil {
		log.Fatalf("Failed to initialize stock subscriber: %v", err)
	}

	if err := stockSubscriber.SubscribeToStockDecreased(); err != nil {
		log.Fatalf("Error subscribing to order.created events: %v", err)
	}

	if err != nil {
		log.Fatalf("failed to connect to product service: %v", err)
	}

	log.Println("Consumer service is running and waiting for events...")

	select {}
}
