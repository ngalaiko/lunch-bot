package rooms

import (
	"time"

	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Room struct {
	ID        ID                `json:"id"`
	Name      string            `json:"name"`
	UserID    users.ID          `json:"userId"`
	Time      time.Time         `json:"time"`
	MemberIDs map[users.ID]bool `json:"memberIds"`
}

func New(userID users.ID, name string) *Room {
	return &Room{
		ID:     ID(uuid.NewString()),
		Name:   name,
		Time:   time.Now(),
		UserID: userID,
		MemberIDs: map[users.ID]bool{
			userID: true,
		},
	}
}
