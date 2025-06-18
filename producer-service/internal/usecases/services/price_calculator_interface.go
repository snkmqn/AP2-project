package services

import "producer-service/internal/core/models"

type PriceCalculator interface {
	CalculateTotalPrice(items []models.OrderItem) float64
}
