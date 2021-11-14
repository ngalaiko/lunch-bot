package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/rolls"
	storage_rolls "lunch/pkg/lunch/rolls/storage"

	"github.com/google/uuid"
)

func migrateRolls(ctx context.Context, from, to storage_rolls.Storage) error {
	toMigrate, err := from.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d rolls to migrate", len(toMigrate))

	migrated, err := from.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated: %w", err)
	}
	log.Printf("[INFO] %d already migrated migrated", len(migrated))

	uniqueKey := func(roll *rolls.Roll) string {
		return fmt.Sprintf("%s:%s:%s", roll.PlaceName, roll.UserID, roll.Time)
	}

	isMigrated := map[string]bool{}
	for _, roll := range migrated {
		isMigrated[uniqueKey(roll)] = true
	}

	for i, roll := range toMigrate {
		if isMigrated[uniqueKey(roll)] {
			log.Printf("[INFO] %d/%d roll already migrated, skipping", i+1, len(toMigrate))
			continue
		}
		if roll.ID == "" {
			roll.ID = uuid.NewString()
		}
		log.Printf("[INFO] %d/%d migrating %+v", i+1, len(toMigrate), roll)
		if err := to.Store(ctx, roll); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", roll, err)
			continue
		}
	}
	return nil
}
