package main

import (
	"context"
	"log"
	"os"

	"lunch/pkg/http"
	"lunch/pkg/jwt"
	storage_jwt_keys "lunch/pkg/jwt/keys/storage"
	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/store"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/iamatypeofwalrus/shim"
)

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
	placesStore   = storage_places.NewDynamoDB(dynamodbStore, "Places")
	boostsStore   = storage_boosts.NewDynamoDB(dynamodbStore, "Boosts")
	rollsStore    = storage_rolls.NewDynamoDB(dynamodbStore, "Rolls")
	jwtKeysStore  = storage_jwt_keys.NewMemory()
)

var (
	roller     = lunch.New(placesStore, boostsStore, rollsStore)
	jwtService = jwt.NewService(jwtKeysStore)
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)
	cfg := &http.Configuration{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}
	handler := shim.New(http.NewHandler(cfg, roller, jwtService), shim.SetDebugLogger(log))
	lambda.StartWithContext(context.Background(), handler.Handle)
}
