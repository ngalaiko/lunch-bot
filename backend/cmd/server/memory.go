// +build !dynamodb

package main

import (
	"log"
	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

func init() {
	log.Println("[INFO] using memory storage")
}

var (
	placesStore  = storage_places.NewMemory()
	boostsStore  = storage_boosts.NewMemory()
	rollsStore   = storage_rolls.NewMemory()
	jwtKeysStore = storage_jwt_keys.NewMemory()
)
