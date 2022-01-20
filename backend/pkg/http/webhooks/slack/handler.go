package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"lunch/pkg/lunch"
	"lunch/pkg/lunch/places"
	"lunch/pkg/users"
	service_users "lunch/pkg/users/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

type Handler struct {
	cfg          *Configuration
	roller       *lunch.Roller
	client       *http.Client
	usersService *service_users.Service
}

func NewHandler(cfg *Configuration, roller *lunch.Roller, usersService *service_users.Service) http.Handler {
	h := &Handler{
		cfg:          cfg,
		roller:       roller,
		client:       &http.Client{},
		usersService: usersService,
	}

	roller.OnRollCreated(h.onRollCreated)
	roller.OnBoostCreated(h.onBoostCreated)
	roller.OnPlaceCreated(h.onPlaceCreated)

	r := chi.NewMux()
	r.With(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded")).Post("/", h.ServeHTTP)

	return r
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	command, actions, challange, err := ParseRequest(r, h.cfg.SigningSecret)
	if err != nil {
		log.Printf("[WARN] failed to parse request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch {
	case actions != nil:
		log.Printf("[INFO] incoming actions: %+v", actions)

		user := &users.User{ID: actions.User.ID, Name: actions.User.Name}
		if err := h.usersService.Create(r.Context(), user); err != nil {
			log.Printf("[ERROR] failed to create user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := users.NewContext(r.Context(), user)
		response := h.handleActions(ctx, actions.ResponseUrl, actions.Actions...)
		if err := respondJSON(w, response); err != nil {
			log.Printf("[ERROR] failed to marshal response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case command != nil:
		log.Printf("[INFO] incoming command: %+v", command)

		user := &users.User{ID: command.UserID, Name: command.UserName}
		if err := h.usersService.Create(r.Context(), user); err != nil {
			log.Printf("[ERROR] failed to create user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := users.NewContext(r.Context(), user)
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
	chances, err := h.roller.ListPlaces(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	bb := []*Block{
		Section(nil, Markdown("*Title*"), Markdown("*Odds*")),
		Divider(),
	}

	for _, chance := range chances {
		bb = append(bb, SectionFields(
			[]*TextBlock{
				PlainText("%s", chance.Name),
				PlainText("%.2f%%", chance.Chance*100),
			},
			WithButton(PlainText("Boost"), "boost", string(chance.ID)),
		))
	}

	return bb, nil
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
	roll, err := h.roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		return Ephemeral(
			fmt.Sprintf("You rolled %s", roll.Place.Name),         // Used in notifications
			Section(Markdown("You rolled *%s*", roll.Place.Name)), // Used in app
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
	if err := h.roller.NewPlace(ctx, placeName); err != nil {
		return InternalServerError(err)
	}
	return Ephemeral(
		fmt.Sprintf("%s added", placeName),
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

func (s *Handler) onBoostCreated(ctx context.Context, boost *lunch.Boost) error {
	place, err := s.roller.GetPlace(ctx, boost.PlaceID)
	if err != nil {
		return fmt.Errorf("failed to get place: %w", err)
	}

	text := fmt.Sprintf("<@%s> boosted %s", boost.UserID, place.Name)
	blocks := Section(Markdown("<@%s> boosted *%s*", boost.UserID, place.Name))

	users, err := s.usersService.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		if user.ID == boost.UserID {
			continue
		}

		user := user
		wg.Go(func() error {
			if err := s.sendMessage(ctx, user, text, blocks); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		})
	}
	return wg.Wait()
}

func (s *Handler) onPlaceCreated(ctx context.Context, place *lunch.Place) error {
	text := fmt.Sprintf("<@%s> added %s", place.AddedBy.ID, place.Name)
	blocks := Section(Markdown("<@%s> added *%s*", place.AddedBy.ID, place.Name))

	users, err := s.usersService.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		if user.ID == place.AddedBy.ID {
			continue
		}

		user := user
		wg.Go(func() error {
			if err := s.sendMessage(ctx, user, text, blocks); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		})
	}
	return wg.Wait()
}

func (s *Handler) onRollCreated(ctx context.Context, roll *lunch.Roll) error {
	place, err := s.roller.GetPlace(ctx, roll.PlaceID)
	if err != nil {
		return fmt.Errorf("failed to get place: %w", err)
	}

	users, err := s.usersService.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	text := fmt.Sprintf("<@%s> rolled %s", roll.UserID, place.Name)
	blocks := Section(Markdown("<@%s> rolled *%s*", roll.UserID, place.Name))

	wg, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		if user.ID == roll.UserID {
			continue
		}
		user := user
		wg.Go(func() error {
			if err := s.sendMessage(ctx, user, text, blocks); err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			return nil
		})
	}
	return wg.Wait()
}

func (s *Handler) sendMessage(ctx context.Context, user *users.User, text string, blocks ...*Block) error {
	type request struct {
		Channel string   `json:"channel"`
		Text    string   `json:"text"`
		Blocks  []*Block `json:"blocks"`
	}

	log.Printf("[INFO] sending message to %s", user.ID)

	type response struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	body, err := json.Marshal(&request{
		Channel: user.ID,
		Text:    text,
		Blocks:  blocks,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.BotAccessToken))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if r.OK {
		return nil
	}

	return fmt.Errorf("failed to send message: %s", r.Error)
}
