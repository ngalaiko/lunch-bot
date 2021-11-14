package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/boosts"
	storage_boosts "lunch/pkg/lunch/boosts/storage"

	"github.com/google/uuid"
)

func migrateBoosts(ctx context.Context, from, to storage_boosts.Storage) error {
	toMigrate, err := from.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d boosts to migrate", len(toMigrate))

	migrated, err := from.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated: %w", err)
	}
	log.Printf("[INFO] %d boosts already migrated", len(migrated))

	uniqueKey := func(boost *boosts.Boost) string {
		return fmt.Sprintf("%s:%s:%s", boost.PlaceName, boost.UserID, boost.Time)
	}

	isMigrated := map[string]bool{}
	for _, boost := range migrated {
		isMigrated[uniqueKey(boost)] = true
	}

	for i, boost := range toMigrate {
		if isMigrated[uniqueKey(boost)] {
			log.Printf("[INFO] %d/%d boost already migrated, skipping", i+1, len(toMigrate))
			continue
		}
		if boost.ID == "" {
			boost.ID = uuid.NewString()
		}
		log.Printf("[INFO] %d/%d migrating %+v", i+1, len(toMigrate), boost)
		if err := to.Store(ctx, boost); err != nil {
			log.Printf("[ERROR] failed to migrate %+v: %s", boost, err)
		}
	}
	return nil
}
