package lunch

type Type uint

const (
	TypeUnknown Type = iota
	TypeRollCreated
	TypeBoostCreated
	TypePlaceCreated
	TypeRoomCreated
	TypeRoomUpdated
)

func (t *Type) String() string {
	switch *t {
	case TypeRoomUpdated:
		return "room_updated"
	case TypeRoomCreated:
		return "room_created"
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

type event struct {
	Type  Type
	Place *Place
	Roll  *Roll
	Boost *Boost
	Room  *Room
}
