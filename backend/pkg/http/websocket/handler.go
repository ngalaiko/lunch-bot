package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"lunch/pkg/lunch"
	"lunch/pkg/lunch/boosts"
	"lunch/pkg/lunch/places"
	"lunch/pkg/lunch/rolls"
	"lunch/pkg/users"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type handler struct {
	roller *lunch.Roller
}

func Handler(roller *lunch.Roller) http.Handler {
	return &handler{
		roller: roller,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	defer conn.Close()

	// TODO: add auth
	ctx := users.NewContext(r.Context(), &users.User{
		ID:   "test",
		Name: "test",
	})

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Printf("[ERROR] failed to read websocket message: %s", err)
			return
		}

		req := &request{}
		if err := json.Unmarshal(msg, req); err != nil {
			log.Printf("[ERROR] failed to unmarshal websocket message: %s", err)

			if err := writeResponse(conn, op, &response{Error: "failed to unmarshal request"}); err != nil {
				log.Printf("[ERROR] failed to write message: %s", err)
				return
			}
			return
		}

		resp, err := h.handle(ctx, req)
		if err != nil {
			log.Printf("[ERROR] failed to handle websocket message: %s", err)

			if err := writeResponse(conn, op, &response{ID: req.ID, Error: "internal error"}); err != nil {
				log.Printf("[ERROR] failed to write message: %s", err)
				return
			}
			return
		}

		if err := writeResponse(conn, op, resp); err != nil {
			log.Printf("[ERROR] failed to write message: %s", err)
			return
		}
	}
}

func (h *handler) handle(ctx context.Context, req *request) (*response, error) {
	switch req.Method {
	case methodAdd:
		return h.handleAdd(ctx, req)
	case methodBoost:
		return h.handleBoost(ctx, req)
	case methodList:
		return h.handleList(ctx, req)
	case methodRoll:
		return h.handleRoll(ctx, req)
	default:
		return &response{ID: req.ID, Error: fmt.Sprintf("unknown method '%s'", req.Method)}, nil
	}
}

func (h *handler) handleAdd(ctx context.Context, req *request) (*response, error) {
	name, ok := req.Params["name"]
	if !ok {
		return &response{ID: req.ID, Error: "'name' parameter must be set"}, nil
	}
	if _, err := h.roller.NewPlace(ctx, name); err != nil {
		return nil, fmt.Errorf("failed to create place: %s", err)
	}
	places, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to list chances: %s", err)
	}
	return &response{ID: req.ID, Places: places}, nil
}

func (h *handler) handleBoost(ctx context.Context, req *request) (*response, error) {
	placeID, ok := req.Params["id"]
	if !ok {
		return &response{ID: req.ID, Error: "'id' parameter must be set"}, nil
	}

	boost, err := h.roller.Boost(ctx, places.ID(placeID), time.Now())
	switch {
	case err == nil:
		return &response{ID: req.ID, Boosts: []*boosts.Boost{boost}}, nil
	case errors.Is(err, lunch.ErrNoPoints):
		return &response{ID: req.ID, Error: "no points left"}, nil
	default:
		return nil, fmt.Errorf("failed to boost: %s", err)
	}
}

func (h *handler) handleList(ctx context.Context, req *request) (*response, error) {
	places, err := h.roller.ListPlaces(ctx, time.Now())
	switch {
	case err == nil:
		return &response{ID: req.ID, Places: places}, nil
	case errors.Is(err, lunch.ErrNoPlaces):
		return &response{ID: req.ID}, nil
	default:
		return nil, fmt.Errorf("failed to list chances: %s", err)
	}
}

func (h *handler) handleRoll(ctx context.Context, req *request) (*response, error) {
	roll, _, err := h.roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		places, err := h.roller.ListPlaces(ctx, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to list chances: %s", err)
		}
		return &response{ID: req.ID, Rolls: []*rolls.Roll{roll}, Places: places}, nil
	case errors.Is(err, lunch.ErrNoPoints):
		return &response{ID: req.ID, Error: "no points left"}, nil
	default:
		return nil, fmt.Errorf("failed to roll: %s", err)
	}
}

func writeResponse(w io.Writer, op ws.OpCode, resp *response) error {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %s", err)
	}
	return wsutil.WriteServerMessage(w, op, bytes)
}
