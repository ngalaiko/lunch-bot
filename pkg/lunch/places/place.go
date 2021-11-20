package places

import (
	"time"

	"lunch/pkg/users"

	"github.com/google/uuid"
)

type Name string

type Place struct {
	ID      string      `dynamodbav:"id"`
	Name    Name        `dynamodbav:"name"`
	AddedAt time.Time   `dynamodbav:"added_at,unixtime"`
	AddedBy *users.User `dynamodbav:"added_by"`
}

func NewPlace(name Name, user *users.User) *Place {
	return &Place{
		ID:      uuid.NewString(),
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
