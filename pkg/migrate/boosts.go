package migrate

import (
	"context"
	"fmt"
	"log"

	storage_boosts "lunch/pkg/lunch/boosts/storage"

	"github.com/google/uuid"
)

func migrateBoosts(ctx context.Context, from, to storage_boosts.Storage) error {
	toMigrate, err := from.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d boosts to migrate", len(toMigrate))

	for i, boost := range toMigrate {
		if boost.ID == "" {
			boost.ID = uuid.NewString()
		}
		log.Printf("[INFO] %d/%d migrating %+v", i, len(toMigrate), boost)
		if err := to.Store(ctx, boost); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", boost, err)
		}
	}
	return nil
}
