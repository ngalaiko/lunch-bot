package places

import (
	"time"

	"lunch/pkg/users"

	"github.com/avelino/slugify"
)

type Place struct {
	ID      string      `json:"-"`
	Name    string      `json:"name"`
	AddedAt time.Time   `json:"added_at"`
	AddedBy *users.User `json:"added_by"`
}

func NewPlace(name string, user *users.User) *Place {
	return &Place{
		ID:      slugify.Slugify(name),
		Name:    name,
		AddedBy: user,
		AddedAt: time.Now(),
	}
}
