package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"lunch/pkg/http/slack"
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/places"
	"lunch/pkg/users"
)

func (s *Server) handleSlack() http.HandlerFunc {
	list := func(ctx context.Context) ([]*slack.Block, error) {
		placeChances, err := s.roller.ListChances(ctx, time.Now())
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

		blocks := []*slack.Block{
			slack.Section(nil, slack.Markdown("*Title*"), slack.Markdown("*Odds*")),
			slack.Divider(),
		}

		for _, p := range pp {
			blocks = append(blocks, slack.SectionFields(
				[]*slack.TextBlock{
					slack.PlainText("%s", p.Name),
					slack.PlainText("%.2f%%", p.Chance),
				},
				slack.WithButton(slack.PlainText("Boost"), "boost", string(p.ID)),
			))
		}

		return blocks, nil
	}

	handleBoost := func(ctx context.Context, w http.ResponseWriter, responseURL string, placeID places.ID) error {
		err := s.roller.Boost(ctx, placeID, time.Now())
		switch {
		case err == nil:
			responseBlocks, err := list(ctx)
			if err != nil {
				return slack.InternalServerError(w, err)
			}
			return slack.ReplaceEphemeralURL(w, responseURL, "Boosting", responseBlocks...)
		case errors.Is(err, lunch.ErrNoPoints):
			return slack.EphemeralURL(w, responseURL, "Failed to boost: no more points left")
		default:
			return slack.InternalServerError(w, err)
		}
	}

	handleSlackActions := func(ctx context.Context, w http.ResponseWriter, responseURL string, actions ...*slack.Action) error {
		if len(actions) != 1 {
			return slack.BadRequest(w, fmt.Errorf("unexpected number of actions: %d", len(actions)))
		}
		action := actions[0]

		log.Printf("[INFO] incoming action: %+v", action)
		switch action.ActionID {
		case "boost":
			return handleBoost(ctx, w, responseURL, places.ID(action.Value))
		default:
			return slack.BadRequest(w, fmt.Errorf("not implemented"))
		}
	}

	handleRoll := func(ctx context.Context, w http.ResponseWriter) error {
		place, err := s.roller.Roll(ctx, time.Now())
		switch {
		case err == nil:
			return slack.InChannel(
				w,
				slack.Text(fmt.Sprintf("Today's lunch place is... %s!", place.Name)),         // Used in notifications
				slack.Section(slack.Markdown("Today's lunch place is... *%s*!", place.Name))) // Used in app
		case errors.Is(err, lunch.ErrNoPoints):
			return slack.Ephemeral(w, "Failed to roll: no more points left")
		case errors.Is(err, lunch.ErrNoPlaces):
			return slack.Ephemeral(w, "No places to choose from, add some!")
		default:
			return slack.InternalServerError(w, err)
		}
	}

	handleAdd := func(ctx context.Context, w http.ResponseWriter, placeName string) error {
		if _, err := s.roller.NewPlace(ctx, placeName); err != nil {
			return slack.InternalServerError(w, err)
		}
		return slack.Ephemeral(
			w,
			slack.Text(fmt.Sprintf("%s added", placeName)),
			slack.Section(slack.Markdown("*%s* added!", placeName)),
		)
	}

	handleList := func(ctx context.Context, w http.ResponseWriter) error {
		responseBlocks, err := list(ctx)
		if err != nil {
			return slack.InternalServerError(w, err)
		}
		return slack.Ephemeral(w, "List", responseBlocks...)
	}

	handleCommand := func(ctx context.Context, w http.ResponseWriter, cmd *slack.CommandRequest) error {
		log.Printf("[INFO] incoming command: %+v", cmd)
		switch cmd.Command {
		case "/roll":
			return handleRoll(ctx, w)
		case "/add":
			return handleAdd(ctx, w, cmd.Text)
		case "/list":
			return handleList(ctx, w)
		default:
			return slack.BadRequest(w, fmt.Errorf("unknown command"))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		command, actions, err := slack.ParseRequest(r)
		if err != nil {
			slack.InternalServerError(w, err)
			return
		}

		if actions != nil {
			ctx := users.NewContext(r.Context(), &users.User{ID: actions.User.ID, Name: actions.User.Name})
			if err := handleSlackActions(ctx, w, actions.ResponseUrl, actions.Actions...); err != nil {
				slack.InternalServerError(w, err)
				return
			}
		} else if command != nil {
			ctx := users.NewContext(r.Context(), &users.User{ID: command.UserID, Name: command.UserName})
			if err := handleCommand(ctx, w, command); err != nil {
				slack.InternalServerError(w, err)
				return
			}
		} else {
			slack.InternalServerError(w, fmt.Errorf("unknown request"))
		}
	}
}
