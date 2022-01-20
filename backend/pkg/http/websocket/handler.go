package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"lunch/pkg/lunch"
	"lunch/pkg/lunch/places"
	"lunch/pkg/users"

	"github.com/go-chi/chi/v5"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

type handler struct {
	roller *lunch.Roller

	openConnections      map[string]io.ReadWriter
	openConnectionsGuard *sync.RWMutex
}

func Handler(roller *lunch.Roller) http.Handler {
	r := chi.NewMux()
	h := &handler{
		roller: roller,

		openConnections:      map[string]io.ReadWriter{},
		openConnectionsGuard: &sync.RWMutex{},
	}
	r.Get("/", h.ServeHTTP)
	roller.OnBoostCreated(h.onBoostCreated)
	roller.OnPlaceCreated(h.onPlaceCreated)
	roller.OnRollCreated(h.onRollCreated)
	return r
}

func (h *handler) onBoostCreated(ctx context.Context, boost *lunch.Boost) error {
	places, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list chances: %s", err)
	}
	return h.broadcast(ws.OpText, &response{Places: places, Boosts: []*lunch.Boost{boost}})
}

func (h *handler) onPlaceCreated(ctx context.Context, place *lunch.Place) error {
	places, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list chances: %s", err)
	}
	return h.broadcast(ws.OpText, &response{Places: places})
}

func (h *handler) onRollCreated(ctx context.Context, roll *lunch.Roll) error {
	pp, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list chances: %s", err)
	}
	return h.broadcast(ws.OpText, &response{Places: pp, Rolls: []*lunch.Roll{roll}})
}

func (h *handler) registerConnection(conn io.ReadWriter) func() {
	id := uuid.NewString()

	h.openConnectionsGuard.Lock()
	h.openConnections[id] = conn
	h.openConnectionsGuard.Unlock()

	return func() {
		h.openConnectionsGuard.Lock()
		delete(h.openConnections, id)
		h.openConnectionsGuard.Unlock()
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := users.FromContext(r.Context()); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	defer h.registerConnection(conn)()
	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil && err != io.EOF {
			if _, isClosedErr := err.(wsutil.ClosedError); isClosedErr {
				return
			}
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

		resp, err := h.handle(r.Context(), req)
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
	case methodPlacesList:
		return h.handlePlacesList(ctx, req)
	case methodPlacesCreate:
		return h.handlePlacesCreate(ctx, req)

	case methodBoostsCreate:
		return h.handleBoostsCreate(ctx, req)
	case methodBoostsList:
		return h.handleBoostsList(ctx, req)

	case methodRollsCreate:
		return h.handleRollsCreate(ctx, req)
	case methodRollsList:
		return h.handleRollsList(ctx, req)
	default:
		return &response{ID: req.ID, Error: fmt.Sprintf("unknown method '%s'", req.Method)}, nil
	}
}

func (h *handler) handlePlacesCreate(ctx context.Context, req *request) (*response, error) {
	name, ok := req.Params["name"]
	if !ok {
		return &response{ID: req.ID, Error: "'name' parameter must be set"}, nil
	}
	if err := h.roller.NewPlace(ctx, name); err != nil {
		return nil, fmt.Errorf("failed to create place: %s", err)
	}
	return &response{ID: req.ID}, nil
}

func (h *handler) handleBoostsCreate(ctx context.Context, req *request) (*response, error) {
	placeID, ok := req.Params["placeId"]
	if !ok {
		return &response{ID: req.ID, Error: "'placeId' parameter must be set"}, nil
	}

	err := h.roller.Boost(ctx, places.ID(placeID), time.Now())
	switch {
	case err == nil:
		return &response{ID: req.ID}, nil
	case errors.Is(err, lunch.ErrNoPoints):
		return &response{ID: req.ID, Error: "no points left"}, nil
	default:
		return nil, fmt.Errorf("failed to boost: %s", err)
	}
}

func (h *handler) handleBoostsList(ctx context.Context, req *request) (*response, error) {
	boosts, err := h.roller.ListBoosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %s", err)
	}
	return &response{ID: req.ID, Boosts: boosts}, nil
}

func (h *handler) handleRollsList(ctx context.Context, req *request) (*response, error) {
	rolls, err := h.roller.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %s", err)
	}
	return &response{ID: req.ID, Rolls: rolls}, nil
}

func (h *handler) handlePlacesList(ctx context.Context, req *request) (*response, error) {
	pp, err := h.roller.ListPlaces(ctx, time.Now())
	switch {
	case err == nil:
		return &response{ID: req.ID, Places: pp}, nil
	case errors.Is(err, lunch.ErrNoPlaces):
		return &response{ID: req.ID}, nil
	default:
		return nil, fmt.Errorf("failed to list chances: %s", err)
	}
}

func (h *handler) handleRollsCreate(ctx context.Context, req *request) (*response, error) {
	roll, err := h.roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		return &response{ID: req.ID, Rolls: []*lunch.Roll{roll}}, nil
	case errors.Is(err, lunch.ErrNoPoints):
		return &response{ID: req.ID, Error: "no points left"}, nil
	default:
		return nil, fmt.Errorf("failed to roll: %s", err)
	}
}

func (h *handler) broadcast(op ws.OpCode, resp *response) error {
	h.openConnectionsGuard.RLock()
	defer h.openConnectionsGuard.RUnlock()

	for _, conn := range h.openConnections {
		if err := writeResponse(conn, op, resp); err != nil {
			log.Printf("[ERROR] failed to write message: %s", err)
		}
	}
	return nil
}

func writeResponse(w io.Writer, op ws.OpCode, resp *response) error {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %s", err)
	}
	return wsutil.WriteServerMessage(w, op, bytes)
}
