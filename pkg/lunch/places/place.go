package places

import (
	"time"

	"lunch/pkg/users"
)

type Name string

type Place struct {
	Name    Name        `dynamodbav:"name"`
	AddedAt time.Time   `dynamodbav:"added_at,unixtime"`
	AddedBy *users.User `dynamodbav:"added_by"`
}

func NewPlace(name Name, user *users.User) *Place {
	return &Place{
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
