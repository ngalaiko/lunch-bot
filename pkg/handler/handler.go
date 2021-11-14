package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	"lunch/pkg/lunch/places"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
	"lunch/pkg/request"
	"lunch/pkg/response"
	"lunch/pkg/store"
	"lunch/pkg/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type actionUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type action struct {
	ActionID string `json:"action_id"`
	Value    string `json:"value"`
}

type actionsRequest struct {
	User        *actionUser `json:"user"`
	Actions     []*action   `json:"actions"`
	ResponseUrl string      `json:"response_url"`
}

type payload struct {
	Payload string `query:"payload"`
}

type commandRequest struct {
	Command  string `query:"command"`
	Text     string `query:"text"`
	UserID   string `query:"user_id"`
	UserName string `query:"user_name"`
}

func mustLoadConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	return cfg
}

var (
	cfg           = mustLoadConfig()
	s3Store       = store.NewS3(cfg)
	dynamodbStore = store.NewDynamoDB(cfg)
	placesStore   = storage_places.NewS3(s3Store)
	boostsStore   = storage_boosts.NewS3(s3Store)
	rollsStore    = storage_rolls.NewS3(s3Store)
	roller        = lunch.New(placesStore, boostsStore, rollsStore)
)

func Handle(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	payload := &payload{}
	if err := request.Parse(req, payload); err == nil && payload.Payload != "" {
		return handlePayload(ctx, payload)
	}

	command := &commandRequest{}
	if err := request.Parse(req, command); err == nil {
		return handleCommand(ctx, command)
	}

	return response.BadRequest(fmt.Errorf("unknown request"))
}

func handlePayload(ctx context.Context, payload *payload) (*events.APIGatewayProxyResponse, error) {
	r := &actionsRequest{}
	if err := json.Unmarshal([]byte(payload.Payload), r); err != nil {
		return response.BadRequest(fmt.Errorf("failed to parse payload"))
	}
	ctx = users.NewContext(ctx, &users.User{
		ID:   r.User.ID,
		Name: r.User.Name,
	})
	return handleActions(ctx, r.ResponseUrl, r.Actions...)
}

func handleActions(ctx context.Context, responseURL string, actions ...*action) (*events.APIGatewayProxyResponse, error) {
	if len(actions) != 1 {
		return response.BadRequest(fmt.Errorf("unexpected number of actions: %d", len(actions)))
	}
	action := actions[0]

	log.Printf("[INFO] incoming action: %+v", action)
	switch action.ActionID {
	case "boost":
		return handleBoost(ctx, responseURL, action.Value)
	default:
		return response.BadRequest(fmt.Errorf("not implemented"))
	}
}

func handleBoost(ctx context.Context, responseURL string, place string) (*events.APIGatewayProxyResponse, error) {
	err := roller.Boost(ctx, place, time.Now())
	switch {
	case err == nil:
		responseBlocks, err := list(ctx)
		if err != nil {
			return response.InternalServerError(err)
		}
		return response.ReplaceEphemeralURL(responseURL, "Boosting", responseBlocks...)
	case errors.Is(err, lunch.ErrNoPoints):
		return response.EphemeralURL(responseURL, "Failed to boost: no more points left")
	default:
		return response.InternalServerError(err)
	}
}

func handleCommand(ctx context.Context, command *commandRequest) (*events.APIGatewayProxyResponse, error) {
	log.Printf("[INFO] incoming command: %+v", command)

	ctx = users.NewContext(ctx, &users.User{
		ID:   command.UserID,
		Name: command.UserName,
	})

	switch command.Command {
	case "/roll":
		return handleRoll(ctx)
	case "/add":
		return handleAdd(ctx, command.Text)
	case "/list":
		return handleList(ctx)
	default:
		return response.BadRequest(fmt.Errorf("unknown command"))
	}
}

func handleRoll(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	place, err := roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		return response.InChannel(
			response.Text(fmt.Sprintf("Today's lunch place is... %s!", place.Name)),            // Used in notifications
			response.Section(response.Markdown("Today's lunch place is... *%s*!", place.Name))) // Used in app
	case errors.Is(err, lunch.ErrNoPoints):
		return response.Ephemeral("Failed to roll: no more points left")
	case errors.Is(err, lunch.ErrNoPlaces):
		return response.Ephemeral("No places to choose from, add some!")
	default:
		return response.InternalServerError(err)
	}
}

func handleAdd(ctx context.Context, place string) (*events.APIGatewayProxyResponse, error) {
	if err := roller.NewPlace(ctx, place); err != nil {
		return response.InternalServerError(err)
	}
	return response.Ephemeral(
		response.Text(fmt.Sprintf("%s added", place)),
		response.Section(response.Markdown("*%s* added!", place)),
	)
}

func handleList(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	responseBlocks, err := list(ctx)
	if err != nil {
		return response.InternalServerError(err)
	}
	return response.Ephemeral("List", responseBlocks...)
}

func list(ctx context.Context) ([]*response.Block, error) {
	placeChances, err := roller.ListChances(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	type placeChance struct {
		Name   places.Name
		Chance float64
	}

	pp := make([]placeChance, 0, len(placeChances))
	for place, chance := range placeChances {
		pp = append(pp, placeChance{
			Name:   place,
			Chance: chance,
		})
	}

	sort.SliceStable(pp, func(i, j int) bool {
		return pp[i].Chance > pp[j].Chance
	})

	blocks := []*response.Block{
		response.Section(nil, response.Markdown("*Title*"), response.Markdown("*Odds*")),
		response.Divider(),
	}

	for _, p := range pp {
		blocks = append(blocks, response.SectionFields(
			[]*response.TextBlock{
				response.PlainText("%s", p.Name),
				response.PlainText("%.2f%%", p.Chance),
			},
			response.WithButton(response.PlainText("Boost"), "boost", string(p.Name)),
		))
	}

	return blocks, nil
}
