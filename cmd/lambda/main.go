package main

import (
	"context"
	"log"
	"os"

	"lunch/pkg/http"
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
	cfg           = mustLoadConfig()
	dynamodbStore = store.NewDynamoDB(cfg)
	placesStore   = storage_places.NewDynamoDB(dynamodbStore)
	boostsStore   = storage_boosts.NewDynamoDB(dynamodbStore)
	rollsStore    = storage_rolls.NewDynamoDB(dynamodbStore)
)

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags)
	cfg := &http.Configuration{}
	if err := cfg.Parse(); err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}
	handler := shim.New(http.NewServer(cfg, boostsStore, placesStore, rollsStore), shim.SetDebugLogger(log))
	lambda.StartWithContext(context.Background(), handler.Handle)
}
