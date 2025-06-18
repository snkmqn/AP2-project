package migrations

import (
	"context"
	"log"
)

type MigrationFunc func(ctx context.Context, db interface{}) error

type Migration struct {
	ID string
	Up MigrationFunc
}

var migrations []Migration

func Register(id string, up MigrationFunc) {
	migrations = append(migrations, Migration{ID: id, Up: up})
}

func Run(ctx context.Context, db interface{}) error {
	for _, m := range migrations {
		log.Printf("Applying migration %s...", m.ID)
		if err := m.Up(ctx, db); err != nil {
			return err
		}
		log.Printf("Migration %s applied successfully.", m.ID)
	}
	return nil
}
