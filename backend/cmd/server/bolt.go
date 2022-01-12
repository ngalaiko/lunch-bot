//go:build !dynamodb
// +build !dynamodb

package main

import (
	"log"

	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/store"
	storage_users "lunch/pkg/users/storage"
)

func init() {
	log.Println("[INFO] using memory storage")
}

var (
	boltStore   = store.MustNewBolt("bolt.db")
	placesStore = storage_places.NewCache(
		storage_places.NewBolt(boltStore),
	)
	boostsStore = storage_boosts.NewCache(
		storage_boosts.NewBolt(boltStore),
	)
	rollsStore = storage_rolls.NewCache(
		storage_rolls.NewBolt(boltStore),
	)
	jwtKeysStore = storage_jwt_keys.NewCache(
		storage_jwt_keys.NewBolt(boltStore),
	)
	usersStore = storage_users.NewCache(
		storage_users.NewBolt(boltStore),
	)
)
