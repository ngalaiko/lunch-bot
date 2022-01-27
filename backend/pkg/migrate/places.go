package migrate

import (
	"context"
	"fmt"
	"log"

	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
)

func migratePlaces(ctx context.Context, from, to storage_places.Storage) error {
	log.Printf("[INFO] migrating places")

	toMigrate, err := from.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to list names to migrate: %w", err)
	}
	log.Printf("[INFO] %d places to migrate", len(toMigrate))

	migrated, err := to.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to list migrated names: %w", err)
	}
	log.Printf("[INFO] %d places already migrated", len(toMigrate))

	isMigrated := map[places.ID]bool{}
	for _, place := range migrated {
		isMigrated[place.ID] = true
	}

	for _, place := range toMigrate {
		if isMigrated[place.ID] {
			continue
		}
		place, err := from.GetByID(ctx, place.ID)
		if err != nil {
			return fmt.Errorf("failed to get item by name: %w", err)
		}
		if err := to.Store(ctx, place); err != nil {
			return fmt.Errorf("failed to store item: %w", err)
		}
	}
	return nil
}
