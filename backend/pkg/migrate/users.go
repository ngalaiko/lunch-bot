package migrate

import (
	"context"
	"fmt"
	"log"

	storage_places "lunch/pkg/lunch/places/storage"
	"lunch/pkg/users"
)

func migrateUsers(ctx context.Context, placesStore storage_places.Storage) error {
	log.Printf("[INFO] migrating users")

	places, err := placesStore.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to list places: %w", err)
	}

	uu := map[string]*users.User{}
	for _, place := range places {
		uu[place.AddedBy.ID] = place.AddedBy
	}

	for _, user := range uu {
		fmt.Printf("\nnikitag: %+v\n\n", user)
	}

	return fmt.Errorf("not implemented")
}
