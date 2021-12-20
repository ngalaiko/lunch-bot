package websocket

import (
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/rolls"
)

type method string

const (
	methodUndefined  method = ""
	methodListPlaces method = "places/list"
	methodListRolls  method = "rolls/list"
	methodCreateRoll method = "rolls/create"
	methodBoost      method = "boost"
	methodAdd        method = "add"
)

type request struct {
	ID     string            `json:"id"`
	Method method            `json:"method"`
	Params map[string]string `json:"params"`
}

type Roll struct {
	*rolls.Roll
	Place *lunch.Place `json:"place"`
}

type response struct {
	ID     string          `json:"id"`
	Places []*lunch.Place  `json:"places,omitempty"`
	Rolls  []*Roll         `json:"rolls,omitempty"`
	Boosts []*boosts.Boost `json:"boosts,omitempty"`
	Error  string          `json:"error,omitempty"`
}
