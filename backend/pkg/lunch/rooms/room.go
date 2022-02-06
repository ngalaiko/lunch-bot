package rooms

import (
	"time"

	"lunch/pkg/users"

	"github.com/google/uuid"
)

type ID string

type Room struct {
	ID        ID
	Name      string
	UserID    users.ID
	Time      time.Time
	MemberIDs map[users.ID]bool
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
