package logger

import (
	"order-service/internal/core/models"
	"order-service/internal/usecases/services"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type NATSSubscriber struct {
	nc           *nats.Conn
	orderHandler services.OrderHandler
}

func NewNATSSubscriber(natsURL string, handler services.OrderHandler) (*NATSSubscriber, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &NATSSubscriber{
		nc:           nc,
		orderHandler: handler,
	}, nil
}

func (s *NATSSubscriber) SubscribeToOrderCreated() error {
	_, err := s.nc.Subscribe("order.created", func(msg *nats.Msg) {
		var order models.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}
		if err := s.orderHandler.Handle(order); err != nil {
			log.Printf("Failed to handle order: %v", err)
		}
	})
	return err
}
