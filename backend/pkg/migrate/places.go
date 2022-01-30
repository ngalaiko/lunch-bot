package migrate

import (
	"context"
	"log"

	storage_places "lunch/pkg/lunch/places/storage"
)

func migratePlaces(ctx context.Context, from, to storage_places.Storage) error {
	log.Printf("[INFO] migrating places")

	return nil
}
