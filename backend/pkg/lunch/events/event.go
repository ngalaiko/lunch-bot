package events

import (
	"fmt"
	"strconv"
	"time"

	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Type string

type Event struct {
	UserID    users.ID     `dynamodbav:"user_id"`
	RoomID    rooms.ID     `dynamodbav:"room_id"`
	Type      Type         `dynamodbav:"type"`
	Timestamp UnixNanoTime `dynamodbav:"timestamp,unixtime"`
	PlaceID   places.ID    `dynamodbav:"place_id"`
	Name      string       `dynamodbav:"name"`
}

type UnixNanoTime time.Time

func (e *UnixNanoTime) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	tv, ok := av.(*types.AttributeValueMemberN)
	if !ok {
		return fmt.Errorf("unexpected type %T for time.Time", av)
	}

	t, err := decodeUnixNanoTime(tv.Value)
	if err != nil {
		return err
	}

	*e = UnixNanoTime(t)
	return nil
}

func decodeUnixNanoTime(n string) (time.Time, error) {
	v, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse %q as unix nano time: %w", n, err)
	}

	return time.Unix(0, v), nil
}
