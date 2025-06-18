package mongo

import (
	"product-service/internal/migrations/migrations"
)

func RegisterMigrations() {
	migrations.Register("001_create_categories", Migration001Categories)
}
