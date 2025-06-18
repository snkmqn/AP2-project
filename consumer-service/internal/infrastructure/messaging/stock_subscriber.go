package messaging

import (
	"consumer-service/internal/core/models"
	"consumer-service/internal/interfaces/repositories"
	"consumer-service/internal/usecases/services"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"proto/generated/ecommerce/inventory"
)

type StockSubscriber struct {
	nc                *nats.Conn
	stockDecrementer  services.StockDecrementer
	productRepository repositories.ProductRepository
	inventoryClient   inventorypb.InventoryServiceClient
}

func NewStockSubscriber(natsURL string, stockDecrementer services.StockDecrementer, productRepository repositories.ProductRepository, grpcServerAddress string) (*StockSubscriber, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := inventorypb.NewInventoryServiceClient(conn)

	return &StockSubscriber{
		nc:                nc,
		stockDecrementer:  stockDecrementer,
		productRepository: productRepository,
		inventoryClient:   client,
	}, nil
}

func (s *StockSubscriber) SubscribeToStockDecreased() error {
	_, err := s.nc.Subscribe("order.created", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))

		var order models.Order

		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}

		if err := s.stockDecrementer.HandleOrder(context.Background(), order); err != nil {
			log.Printf("Failed to handle stock decrease: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to subscribe to NATS subject: %v", err)
	}
	return err
}
