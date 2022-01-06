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
	"lunch/pkg/users"

	"github.com/go-chi/chi/v5"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type handler struct {
	roller *lunch.Roller
}

func Handler(roller *lunch.Roller) http.Handler {
	r := chi.NewMux()
	r.Get("/", (&handler{roller: roller}).ServeHTTP)
	return r
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
	defer conn.Close()

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
	if _, err := h.roller.NewPlace(ctx, name); err != nil {
		return nil, fmt.Errorf("failed to create place: %s", err)
	}
	places, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to list chances: %s", err)
	}
	return &response{ID: req.ID, Places: places}, nil
}

func (h *handler) handleBoostsCreate(ctx context.Context, req *request) (*response, error) {
	placeID, ok := req.Params["placeId"]
	if !ok {
		return &response{ID: req.ID, Error: "'placeId' parameter must be set"}, nil
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

func (h *handler) handleRollsList(ctx context.Context, req *request) (*response, error) {
	rr, err := h.roller.ListRolls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rolls: %s", err)
	}

	pp, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil && !errors.Is(err, lunch.ErrNoPlaces) {
		return nil, fmt.Errorf("failed to list places: %s", err)
	}

	placesByID := make(map[places.ID]*lunch.Place, len(pp))
	for _, p := range pp {
		placesByID[p.ID] = p
	}

	rolls := make([]*Roll, 0, len(rr))
	for _, r := range rr {
		rolls = append(rolls, &Roll{
			Roll:  r,
			Place: placesByID[r.PlaceID],
		})
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
	roll, _, err := h.roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		pp, err := h.roller.ListPlaces(ctx, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to list chances: %s", err)
		}
		placesByID := make(map[places.ID]*lunch.Place)
		for _, place := range pp {
			placesByID[place.ID] = place
		}
		return &response{ID: req.ID, Rolls: []*Roll{
			{
				Roll:  roll,
				Place: placesByID[roll.PlaceID],
			},
		}}, nil
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
