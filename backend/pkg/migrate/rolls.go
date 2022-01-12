package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/rolls"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

func migrateRolls(ctx context.Context, from, to storage_rolls.Storage) error {
	log.Printf("[INFO] migrating rolls")

	toMigrate, err := from.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d rolls to migrate", len(toMigrate))

	migrated, err := to.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated: %w", err)
	}
	log.Printf("[INFO] %d already migrated migrated", len(migrated))

	isMigrated := map[rolls.ID]bool{}
	for _, roll := range migrated {
		isMigrated[roll.ID] = true
	}

	for i, roll := range toMigrate {
		if isMigrated[roll.ID] {
			log.Printf("[INFO] %d/%d roll already migrated, skipping", i+1, len(toMigrate))
			continue
		}
		log.Printf("[INFO] %d/%d migrating %+v", i+1, len(toMigrate), roll)
		if err := to.Store(ctx, roll); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", roll, err)
			continue
		}
	}
	return nil
}
