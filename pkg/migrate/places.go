package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
)

func migratePlaces(ctx context.Context, from, to storage_places.Storage) error {
	toMigrate, err := from.ListNames(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d places to migrate", len(toMigrate))

	migrated, err := to.ListNames(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated names: %w", err)
	}
	log.Printf("[INFO] %d places already migrated", len(toMigrate))

	isMigrated := map[places.Name]bool{}
	for _, name := range migrated {
		isMigrated[name] = true
	}

	for i, name := range toMigrate {
		if isMigrated[name] {
			log.Printf("[INFO] %d/%d %+v migrated, skipping", i, len(toMigrate), name)
			continue
		}
		place, err := from.GetByName(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to get item by name: %w", err)
		}
		log.Printf("[INFO] %d/%d migrating %+v", i, len(toMigrate), place)
		if err := to.Store(ctx, place); err != nil {
			return fmt.Errorf("failed to store item: %w", err)
		}
	}
	return nil
}
