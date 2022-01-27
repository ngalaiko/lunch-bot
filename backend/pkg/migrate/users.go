package migrate

import (
	"context"
	"log"

	storage_places "lunch/pkg/lunch/places/storage"
	storage_users "lunch/pkg/users/storage"
)

func migrateUsers(ctx context.Context, placesStore storage_places.Storage, usersStorage storage_users.Storage) error {
	log.Printf("[INFO] migrating users")

	return nil
}
