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
	placesStore   = storage_places.NewDynamoDB(dynamodbStore, "lunch-production-webapp-places")
	boostsStore   = storage_boosts.NewDynamoDB(dynamodbStore, "lunch-production-webapp-boosts")
	rollsStore    = storage_rolls.NewDynamoDB(dynamodbStore, "lunch-production-webapp-rolls")
	jwtKeysStore  = storage_jwt_keys.NewMemory() // TODO: use DynamoDB
)
