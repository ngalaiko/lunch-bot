package migrate

import (
	"context"
	"fmt"

	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_boosts_v2 "lunch/pkg/lunch/boosts/storage/v2"
	"lunch/pkg/lunch/events"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_places_v2 "lunch/pkg/lunch/places/storage/v2"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	storage_rolls_v2 "lunch/pkg/lunch/rolls/storage/v2"
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

	placesv1 = storage_places.NewDynamoDB(dynamodbStore, "lunch-production-webapp-places")
	boostsv1 = storage_boosts.NewDynamoDB(dynamodbStore, "lunch-production-webapp-boosts")
	rollsv1  = storage_rolls.NewDynamoDB(dynamodbStore, "lunch-production-webapp-rolls")

	eventsStore = events.NewDynamoDBStore(dynamodbStore, "lunch-production-webapp-events")
	placesv2    = storage_places_v2.New(eventsStore)
	boostsv2    = storage_boosts_v2.New(eventsStore)
	rollsv2     = storage_rolls_v2.New(eventsStore)

	sturdyRoomID = rooms.ID("69c83096-995a-48ce-b843-80a926b0a9ec")
)

func Run(ctx context.Context) error {
	places, err := placesv1.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load places: %w", err)
	}
	for _, place := range places {
		place.RoomID = sturdyRoomID

		fmt.Printf("place.ID: %v %d\n", place.ID, place.Time.UnixNano())
		if err := placesv2.Create(ctx, place); err != nil {
			return fmt.Errorf("failed to migrate place: %w", err)
		}
	}

	boosts, err := boostsv1.ListBoosts(ctx)
	if err != nil {
		return fmt.Errorf("failed to load boosts: %w", err)
	}
	for _, boost := range boosts {
		boost.RoomID = sturdyRoomID
		if err := boostsv2.Create(ctx, boost); err != nil {
			return fmt.Errorf("failed to migrate boost: %w", err)
		}
		fmt.Printf("boost.ID: %v\n", boost.ID)
	}

	rolls, err := rollsv1.ListRolls(ctx)
	if err != nil {
		return fmt.Errorf("failed to load rolls: %w", err)
	}
	for _, roll := range rolls {
		roll.RoomID = sturdyRoomID
		if err := rollsv2.Create(ctx, roll); err != nil {
			return fmt.Errorf("failed to migrate roll: %w", err)
		}
		fmt.Printf("roll.ID: %v\n", roll.ID)
	}
	return nil
}
