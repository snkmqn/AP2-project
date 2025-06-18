package validators

import "api-gateway/internal/core/models"

type UserValidator interface {
	Validate(user models.User) error
}
