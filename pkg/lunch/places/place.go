package places

import (
	"time"

	"lunch/pkg/users"
)

type Place struct {
	Name    string      `json:"name"`
	AddedAt time.Time   `json:"added_at"`
	AddedBy *users.User `json:"added_by"`
}

func NewPlace(name string, user *users.User) *Place {
	return &Place{
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
