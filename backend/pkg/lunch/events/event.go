package events

import (
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
)

type Type uint

const (
	TypeUnknown Type = iota
	TypeRollCreated
	TypeBoostCreated
	TypePlaceCreated
)

func (t *Type) String() string {
	switch *t {
	case TypeRollCreated:
		return "roll_created"
	case TypeBoostCreated:
		return "boost_created"
	case TypePlaceCreated:
		return "place_created"
	default:
		return "unknown"
	}
}

type Event struct {
	Type  Type
	Place *places.Place
	Roll  *rolls.Roll
	Boost *boosts.Boost
}
