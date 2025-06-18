package services

import (
	"product-service/internal/core/models"
	"fmt"
)

type OrderHandler interface {
	Handle(order models.Order) error
}

type DefaultOrderHandler struct{}

func NewOrderHandler() *DefaultOrderHandler {
	return &DefaultOrderHandler{}
}

func (h *DefaultOrderHandler) Handle(order models.Order) error {
	fmt.Printf("Received order created event: Order ID: %s, User ID: %s, Status: %s\n", order.OrderID, order.UserID, order.Status)
	return nil
}
