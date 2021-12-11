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
	placesStorage = storage_places.NewDynamoDB(dynamodbStore)
	boostsStore   = storage_boosts.NewDynamoDB(dynamodbStore)
	rollsStore    = storage_rolls.NewDynamoDB(dynamodbStore)
)

func Run(ctx context.Context) error {
	return fmt.Errorf("nothing to do")
}
