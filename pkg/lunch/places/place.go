package places

import (
	"time"

	"lunch/pkg/users"
)

type Name string

type Place struct {
	Name    Name        `dynamodbav:"name" json:"name"`
	AddedAt time.Time   `dynamodbav:"added_at,unixtime" json:"added_at"`
	AddedBy *users.User `dynamodbav:"added_by" json:"added_by"`
}

func NewPlace(name Name, user *users.User) *Place {
	return &Place{
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
