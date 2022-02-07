package migrate

import (
	"context"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/events"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/lunch/rooms"
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

	eventsStore = events.NewDynamoDBStore(dynamodbStore, "lunch-production-webapp-events")
	places      = storage_places.New(eventsStore)
	boosts      = storage_boosts.New(eventsStore)
	rolls       = storage_rolls.New(eventsStore)

	sturdyRoomID = rooms.ID("69c83096-995a-48ce-b843-80a926b0a9ec")
)

func Run(ctx context.Context) error {
	return nil
}
