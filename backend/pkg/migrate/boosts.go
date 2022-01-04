package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
)

func migrateBoosts(ctx context.Context, from, to storage_boosts.Storage) error {
	toMigrate, err := from.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d boosts to migrate", len(toMigrate))

	migrated, err := to.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated: %w", err)
	}
	log.Printf("[INFO] %d boosts already migrated", len(migrated))

	isMigrated := map[boosts.ID]bool{}
	for _, boost := range migrated {
		isMigrated[boost.ID] = true
	}

	for i, boost := range toMigrate {
		if isMigrated[boost.ID] {
			log.Printf("[INFO] %d/%d boost already migrated, skipping", i+1, len(toMigrate))
			continue
		}
		log.Printf("[INFO] %d/%d migrating %+v", i+1, len(toMigrate), boost)
		if err := to.Store(ctx, boost); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", boost, err)
		}
	}
	return nil
}
