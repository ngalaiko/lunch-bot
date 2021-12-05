package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"lunch/pkg/lunch"
	"lunch/pkg/lunch/places"
	"lunch/pkg/users"
)

type Handler struct {
	roller *lunch.Roller
	client *http.Client
}

func NewHandler(roller *lunch.Roller) *Handler {
	return &Handler{
		roller: roller,
		client: &http.Client{},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	command, actions, challange, err := ParseRequest(r)
	if err != nil {
		log.Printf("[ERROR] failed to parse request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch {
	case actions != nil:
		log.Printf("[INFO] incoming actions: %+v", actions)
		ctx := users.NewContext(r.Context(), &users.User{ID: actions.User.ID, Name: actions.User.Name})
		response := h.handleActions(ctx, actions.ResponseUrl, actions.Actions...)
		if err := respondJSON(w, response); err != nil {
			log.Printf("[ERROR] failed to marshal response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case command != nil:
		log.Printf("[INFO] incoming command: %+v", command)
		ctx := users.NewContext(r.Context(), &users.User{ID: command.UserID, Name: command.UserName})
		response := h.handleCommand(ctx, command)
		if err := respondJSON(w, response); err != nil {
			log.Printf("[ERROR] failed to marshal response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case challange != nil:
		log.Printf("[INFO] incoming challange: %+v", challange)
		if err := respondPlainText(w, []byte(challange.Challenge)); err != nil {
			log.Printf("[ERROR] failed to write challange: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	default:
		log.Printf("[ERROR] unknown request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func respondPlainText(w http.ResponseWriter, body []byte) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write(body)
	return err
}

func respondJSON(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(body)
}

func (h *Handler) list(ctx context.Context) ([]*Block, error) {
	placeChances, err := h.roller.ListChances(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	type placeChance struct {
		ID     places.ID
		Name   places.Name
		Chance float64
	}

	pp := make([]placeChance, 0, len(placeChances))
	for place, chance := range placeChances {
		pp = append(pp, placeChance{
			ID:     place.ID,
			Name:   place.Name,
			Chance: chance,
		})
	}

	sort.Slice(pp, func(i, j int) bool {
		return pp[i].Name < pp[j].Name
	})

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Chance > pp[j].Chance
	})

	blocks := []*Block{
		Section(nil, Markdown("*Title*"), Markdown("*Odds*")),
		Divider(),
	}

	for _, p := range pp {
		blocks = append(blocks, SectionFields(
			[]*TextBlock{
				PlainText("%s", p.Name),
				PlainText("%.2f%%", p.Chance),
			},
			WithButton(PlainText("Boost"), "boost", string(p.ID)),
		))
	}

	return blocks, nil
}

func (h *Handler) asyncPost(url string, msg *Message) error {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(msg); err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	resp, err := h.client.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("failed to post response to slack: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post response to slack: %s", resp.Status)
	}

	return nil
}

func (h *Handler) handleBoost(ctx context.Context, responseURL string, placeID places.ID) error {
	err := h.roller.Boost(ctx, placeID, time.Now())
	switch {
	case err == nil:
		responseBlocks, err := h.list(ctx)
		if err != nil {
			return h.asyncPost(responseURL, InternalServerError(err))
		}
		return h.asyncPost(responseURL, ReplaceEphemeral("Boosting", responseBlocks...))
	case errors.Is(err, lunch.ErrNoPoints):
		return h.asyncPost(responseURL, Ephemeral("Failed to boost: no more points left"))
	default:
		return h.asyncPost(responseURL, InternalServerError(err))
	}
}

func (h *Handler) handleActions(ctx context.Context, responseURL string, actions ...*Action) error {
	if len(actions) != 1 {
		return h.asyncPost(responseURL, BadRequest(fmt.Errorf("unexpected number of actions: %d", len(actions))))
	}
	action := actions[0]

	log.Printf("[INFO] incoming action: %+v", action)
	switch action.ActionID {
	case "boost":
		if err := h.handleBoost(ctx, responseURL, places.ID(action.Value)); err != nil {
			return h.asyncPost(responseURL, InternalServerError(err))
		}
		return nil
	default:
		return h.asyncPost(responseURL, BadRequest(fmt.Errorf("unknown action '%s'", action.ActionID)))
	}
}

func (h *Handler) handleRoll(ctx context.Context) *Message {
	place, err := h.roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		return InChannel(
			Text(fmt.Sprintf("Today's lunch place is... %s!", place.Name)),   // Used in notifications
			Section(Markdown("Today's lunch place is... *%s*!", place.Name)), // Used in app
		)
	case errors.Is(err, lunch.ErrNoPoints):
		return Ephemeral("Failed to roll: no more points left")
	case errors.Is(err, lunch.ErrNoPlaces):
		return Ephemeral("No places to choose from, add some!")
	default:
		return InternalServerError(err)
	}
}

func (h *Handler) handleAdd(ctx context.Context, placeName string) *Message {
	if _, err := h.roller.NewPlace(ctx, placeName); err != nil {
		return InternalServerError(err)
	}
	return Ephemeral(
		Text(fmt.Sprintf("%s added", placeName)),
		Section(Markdown("*%s* added!", placeName)),
	)
}

func (h *Handler) handleList(ctx context.Context) *Message {
	responseBlocks, err := h.list(ctx)
	if err != nil {
		return InternalServerError(err)
	}
	return Ephemeral("List", responseBlocks...)
}

func (h *Handler) handleCommand(ctx context.Context, cmd *CommandRequest) *Message {
	log.Printf("[INFO] incoming command: %+v", cmd)
	switch cmd.Command {
	case "/roll":
		return h.handleRoll(ctx)
	case "/add":
		return h.handleAdd(ctx, cmd.Text)
	case "/list":
		return h.handleList(ctx)
	default:
		return BadRequest(fmt.Errorf("unknown command '%s'", cmd.Command))
	}
}
