package events

import (
	"context"
	"fmt"
	"time"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/store"
	"lunch/pkg/users"
)

type dynamoDB struct {
	db        *store.DynamoDB
	tableName string
}

func NewDynamoDBStore(db *store.DynamoDB, tableName string) *dynamoDB {
	return &dynamoDB{
		db:        db,
		tableName: tableName,
	}
}

func (d *dynamoDB) Create(ctx context.Context, event *Event) error {
	if err := d.db.Execute(ctx, fmt.Sprintf(`
		INSERT INTO "%s" value {
			'user_id': ?,
			'room_id': ?,
			'type': ?,
			'timestamp': ?,
			'place_id': ?,
			'name': ?
		}
	`, d.tableName), event.UserID, event.RoomID, event.Type, time.Time(event.Timestamp).UnixNano(), event.PlaceID, event.Name); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *dynamoDB) ByUserID(ctx context.Context, userID users.ID, types ...Type) ([]*Event, error) {
	ee := []*Event{}
	if err := d.db.Query(ctx, &ee, fmt.Sprintf(`
		SELECT * FROM "%s"
		WHERE 'user_id' = ?
	`, d.tableName), userID); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if len(types) == 0 {
		return ee, nil
	}

	includeType := make(map[Type]bool)
	for _, t := range types {
		includeType[t] = true
	}
	filteredEvents := make([]*Event, 0, len(ee))
	for _, e := range ee {
		if includeType[e.Type] {
			filteredEvents = append(filteredEvents, e)
		}
	}
	return filteredEvents, nil
}

func (d *dynamoDB) ByRoomID(ctx context.Context, roomID rooms.ID, types ...Type) ([]*Event, error) {
	ee := []*Event{}
	if err := d.db.Query(ctx, &ee, fmt.Sprintf(`
		SELECT * FROM "%s"."room_id.timestamp"
		WHERE room_id = ?
	`, d.tableName), roomID); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if len(types) == 0 {
		return ee, nil
	}

	includeType := make(map[Type]bool)
	for _, t := range types {
		includeType[t] = true
	}
	filteredEvents := make([]*Event, 0, len(ee))
	for _, e := range ee {
		if includeType[e.Type] {
			filteredEvents = append(filteredEvents, e)
		}
	}
	return filteredEvents, nil
}
