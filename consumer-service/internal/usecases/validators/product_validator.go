package validators

import (
	"consumer-service/internal/core/models"
	"consumer-service/internal/errors"
	"consumer-service/internal/infrastructure/utils/uuid"
)

func ValidateProductForCreation(product models.Product) error {
	if product.Name == "" {
		return errors.ErrMissingName
	}
	if product.Price <= 0 {
		return errors.ErrInvalidPrice
	}
	if product.Stock <= 0 {
		return errors.ErrInvalidStock
	}
	if product.Description == "" {
		return errors.ErrMissingDescription
	}
	if !uuid.IsValidUUID(product.CategoryID) {
		return errors.ErrInvalidCategoryID
	}
	return nil
}

func ValidateProductForUpdate(product models.Product) error {
	if product.Name == "" {
		return errors.ErrMissingName
	}
	if product.Price <= 0 {
		return errors.ErrInvalidPrice
	}
	if product.Stock < 0 {
		return errors.ErrInvalidStock
	}
	if !uuid.IsValidUUID(product.CategoryID) {
		return errors.ErrInvalidCategoryID
	}
	return nil
}
