package migrate

import (
	"context"
	"fmt"
	"log"

	storage_rolls "lunch/pkg/lunch/rolls/storage"

	"github.com/google/uuid"
)

func migrateRolls(ctx context.Context, from, to storage_rolls.Storage) error {
	toMigrate, err := from.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d rolls to migrate", len(toMigrate))

	for i, roll := range toMigrate {
		if roll.ID == "" {
			roll.ID = uuid.NewString()
		}
		log.Printf("[INFO] %d/%d migrating %+v", i, len(toMigrate), roll)
		if err := to.Store(ctx, roll); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", roll, err)
			continue
		}
	}
	return nil
}
