package repositories

import (
	"consumer-service/internal/core/models"
	"context"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id string) (*models.Order, error)
	UpdateOrder(ctx context.Context, id string, status string) error
	GetOrdersByUserID (ctx context.Context, userID string) ([]*models.Order, error)
	DeleteOrdersByUserID (ctx context.Context, userID string) error
}
