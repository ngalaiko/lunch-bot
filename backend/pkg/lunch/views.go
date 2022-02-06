package lunch

import (
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/lunch/rooms"
	"lunch/pkg/users"
)

type Place struct {
	*places.Place
	User   *users.User `json:"user"`
	Chance float64     `json:"chance"`
}

type Boost struct {
	*boosts.Boost
	User  *users.User   `json:"user"`
	Place *places.Place `json:"place"`
}

type Roll struct {
	*rolls.Roll
	User  *users.User   `json:"user"`
	Place *places.Place `json:"place"`
}

type Room struct {
	*rooms.Room
	User    *users.User   `json:"user"`
	Members []*users.User `json:"members"`
}
