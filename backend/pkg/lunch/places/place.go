package places

import (
	"time"

	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Place struct {
	ID        ID        `json:"id"`
	Name      string    `json:"name"`
	Time      time.Time `json:"time"`
	UserID    users.ID  `json:"userId"`
	RoomID    rooms.ID  `json:"roomId"`
	IsDeleted bool      `json:"-"`
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
