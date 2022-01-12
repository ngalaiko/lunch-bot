//go:build dynamodb
// +build dynamodb

package main

import (
	"context"
	"log"

	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/store"
	storage_users "lunch/pkg/users/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func init() {
	log.Println("[INFO] using dynamodb storage")
}

func mustLoadConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	return cfg
}

var (
	awsConfig     = mustLoadConfig()
	dynamodbStore = store.NewDynamoDB(awsConfig)
	placesStore   = storage_places.NewCache(
		storage_places.NewDynamoDB(dynamodbStore, "lunch-production-webapp-places"),
	)
	boostsStore = storage_boosts.NewCache(
		storage_boosts.NewDynamoDB(dynamodbStore, "lunch-production-webapp-boosts"),
	)
	rollsStore = storage_rolls.NewCache(
		storage_rolls.NewDynamoDB(dynamodbStore, "lunch-production-webapp-rolls"),
	)
	jwtKeysStore = storage_jwt_keys.NewCache(
		storage_jwt_keys.NewDynamoDB(dynamodbStore, "lunch-production-webapp-private-keys"),
	)
	usersStore = storage_users.NewCache(
		storage_users.NewDynamoDB(dynamodbStore, "lunch-production-webapp-users"),
	)
)
