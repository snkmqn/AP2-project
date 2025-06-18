package services

import "api-gateway/internal/core/models"

type PriceCalculator interface {
	CalculateTotalPrice(items []models.OrderItem) float64
}
