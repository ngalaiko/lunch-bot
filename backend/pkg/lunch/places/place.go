package places

import (
	"time"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Place struct {
	ID     ID        `dynamodbav:"id" json:"id"`
	Name   string    `dynamodbav:"name" json:"name"`
	Time   time.Time `dynamodbav:"added_at,unixtime" json:"time"`
	UserID users.ID  `dynamodbav:"user_id" json:"userId"`
	RoomID rooms.ID  `json:"roomId"`
}

func NewPlace(roomID rooms.ID, userID users.ID, name string) *Place {
	return &Place{
		ID:     ID(uuid.NewString()),
		Name:   name,
		Time:   time.Now(),
		UserID: userID,
		RoomID: roomID,
	}
}
