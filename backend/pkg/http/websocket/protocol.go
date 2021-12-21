package websocket

import (
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/rolls"
)

type method string

const (
	methodUndefined    method = ""
	methodPlacesList   method = "places/list"
	methodPlacesCreate method = "places/create"
	methodRollsList    method = "rolls/list"
	methodRollsCreate  method = "rolls/create"
	methodBoostsCreate method = "boosts/create"
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
