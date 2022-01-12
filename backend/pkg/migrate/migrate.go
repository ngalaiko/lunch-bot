package migrate

import (
	"context"
	"fmt"

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
	dynamodbStore = store.NewDynamoDB(cfg)
)

func Run(ctx context.Context) error {
	if err := migrateUsers(ctx, storage_places.NewDynamoDB(dynamodbStore, "lunch-production-webapp-places")); err != nil {
		return fmt.Errorf("failed to migrate users: %w", err)
	}

	if err := migratePlaces(
		ctx,
		storage_places.NewDynamoDB(dynamodbStore, "Places"),
		storage_places.NewDynamoDB(dynamodbStore, "lunch-production-webapp-places"),
	); err != nil {
		return fmt.Errorf("failed to migrate places: %w", err)
	}

	if err := migrateBoosts(
		ctx,
		storage_boosts.NewDynamoDB(dynamodbStore, "Boosts"),
		storage_boosts.NewDynamoDB(dynamodbStore, "lunch-production-webapp-boosts"),
	); err != nil {
		return fmt.Errorf("failed to migrate boosts: %w", err)
	}

	if err := migrateRolls(
		ctx,
		storage_rolls.NewDynamoDB(dynamodbStore, "Rolls"),
		storage_rolls.NewDynamoDB(dynamodbStore, "lunch-production-webapp-rolls"),
	); err != nil {
		return fmt.Errorf("failed to migrate boosts: %w", err)
	}

	return nil
}
