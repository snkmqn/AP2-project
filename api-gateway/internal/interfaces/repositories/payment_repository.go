package repositories

import (
	"context"
	"api-gateway/internal/core/models"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*models.Payment, error)
	DeletePayment(ctx context.Context, id string) error
}
