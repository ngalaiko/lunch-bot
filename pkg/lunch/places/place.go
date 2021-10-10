package places

import (
	"time"

	"lunch/pkg/users"
)

type Name string

type Place struct {
	Name    Name        `json:"name"`
	AddedAt time.Time   `json:"added_at"`
	AddedBy *users.User `json:"added_by"`
}

func NewPlace(name Name, user *users.User) *Place {
	return &Place{
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
