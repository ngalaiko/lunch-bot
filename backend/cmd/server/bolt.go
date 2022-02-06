//go:build !dynamodb
// +build !dynamodb

package main

import (
	"log"

	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	"lunch/pkg/lunch/events"
	"lunch/pkg/store"
	storage_users "lunch/pkg/users/storage"
)

func init() {
	log.Println("[INFO] using memory storage")
}

var (
	boltStore    = store.MustNewBolt("bolt.db")
	jwtKeysStore = storage_jwt_keys.NewCache(
		storage_jwt_keys.NewBolt(boltStore),
	)
	usersStore = storage_users.NewCache(
		storage_users.NewBolt(boltStore),
	)
	eventsStorage = events.NewCache(
		events.NewBoltStorage(boltStore),
	)
)
