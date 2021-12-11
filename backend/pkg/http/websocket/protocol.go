package websocket

import (
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/rolls"
)

type method string

const (
	methodUndefined method = ""
	methodRoll      method = "roll"
	methodList      method = "list"
	methodBoost     method = "boost"
	methodAdd       method = "add"
)

type request struct {
	ID     string            `json:"id"`
	Method method            `json:"method"`
	Params map[string]string `json:"params"`
}

type response struct {
	ID     string          `json:"id"`
	Places []*lunch.Place  `json:"places,omitempty"`
	Rolls  []*rolls.Roll   `json:"rolls,omitempty"`
	Boosts []*boosts.Boost `json:"boosts,omitempty"`
	Error  string          `json:"error,omitempty"`
}
