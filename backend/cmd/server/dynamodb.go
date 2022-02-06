//go:build dynamodb
// +build dynamodb

package main

import (
	"context"
	"log"

	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	"lunch/pkg/lunch/events"
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

	jwtKeysStore = storage_jwt_keys.NewCache(
		storage_jwt_keys.NewDynamoDB(dynamodbStore, "lunch-production-webapp-private-keys"),
	)
	usersStore = storage_users.NewCache(
		storage_users.NewDynamoDB(dynamodbStore, "lunch-production-webapp-users"),
	)
	eventsStorage = events.NewCache(
		events.NewDynamoDBStore(dynamodbStore, "lunch-production-webapp-events"),
	)
)
