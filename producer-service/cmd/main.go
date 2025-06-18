package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"producer-service/internal/infrastructure/messaging"
	"producer-service/internal/usecases/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	natsURL := os.Getenv("NATS_URL")

	handler := services.NewOrderHandler()

	subscriber, err := messaging.NewNATSSubscriber(natsURL, handler)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	err = subscriber.SubscribeToOrderCreated()
	if err != nil {
		log.Fatalf("Error subscribing to order.created: %v", err)
	}

	log.Println("Waiting for order.created events...")
	select {}
}
