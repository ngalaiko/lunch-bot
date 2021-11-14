package migrate

import (
	"context"
	"fmt"
	"log"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/store"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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
	s3Store       = store.NewS3(cfg)
	dynamodbStore = store.NewDynamoDB(cfg)
)

func Run(ctx context.Context) error {
	log.Printf("[INFO] migrating places")
	if err := migratePlaces(ctx, storage_places.NewS3(s3Store), storage_places.NewDynamoDB(dynamodbStore)); err != nil {
		return fmt.Errorf("failed to migrate places: %w", err)
	}

	log.Printf("[INFO] migrating boosts")
	if err := migrateBoosts(ctx, storage_boosts.NewS3(s3Store), storage_boosts.NewDynamoDB(dynamodbStore)); err != nil {
		return fmt.Errorf("failed to migrate boosts: %w", err)
	}

	log.Printf("[INFO] migrating rolls")
	if err := migrateRolls(ctx, storage_rolls.NewS3(s3Store), storage_rolls.NewDynamoDB(dynamodbStore)); err != nil {
		return fmt.Errorf("failed to migrate boosts: %w", err)
	}

	return nil
}
