package websocket

import (
	"lunch/pkg/lunch"
)

type method string

const (
	methodUndefined    method = ""
	methodPlacesList   method = "places/list"
	methodPlacesCreate method = "places/create"
	methodRollsList    method = "rolls/list"
	methodRollsCreate  method = "rolls/create"
	methodBoostsCreate method = "boosts/create"
	methodBoostsList   method = "boosts/list"
	methodRoomsList    method = "rooms/list"
	methodRoomsCreate  method = "rooms/create"
)

type request struct {
	ID     string            `json:"id"`
	Method method            `json:"method"`
	Params map[string]string `json:"params"`
}

type response struct {
	ID     string         `json:"id,omitempty"`
	Places []*lunch.Place `json:"places,omitempty"`
	Rolls  []*lunch.Roll  `json:"rolls,omitempty"`
	Boosts []*lunch.Boost `json:"boosts,omitempty"`
	Rooms  []*lunch.Room  `json:"rooms,omitempty"`
	Error  string         `json:"error,omitempty"`
}
