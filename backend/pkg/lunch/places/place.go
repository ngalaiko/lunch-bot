package places

import (
	"time"

	"lunch/pkg/users"

	"github.com/google/uuid"
)

type Name string

type ID string

type Place struct {
	ID      ID          `dynamodbav:"id" json:"id"`
	Name    Name        `dynamodbav:"name" json:"name"`
	AddedAt time.Time   `dynamodbav:"added_at,unixtime" json:"addedAt"`
	AddedBy *users.User `dynamodbav:"added_by" json:"-"`
}

func NewPlace(name Name, user *users.User) *Place {
	return &Place{
		ID:      ID(uuid.NewString()),
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
