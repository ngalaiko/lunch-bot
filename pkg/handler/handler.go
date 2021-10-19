package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"lunch/pkg/lunch"
	"lunch/pkg/request"
	"lunch/pkg/response"
	"lunch/pkg/store"
	"lunch/pkg/users"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type CommandRequest struct {
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
	cfg     = mustLoadConfig()
	s3Store = store.NewS3(cfg)
	roller  = lunch.New(s3Store)
)

func Handle(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	command := &CommandRequest{}
	if err := request.Parse(req, &command); err != nil {
		return response.BadRequest(err)
	}

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
	case "/boost":
		return handeBoost(ctx)
	default:
		return response.BadRequest(fmt.Errorf("unknown command"))
	}
}

func handleRoll(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	place, err := roller.Roll(ctx, time.Now())
	switch {
	case err == nil:
		return response.InChannel(response.Section(response.Markdown("Today's lunch place is... *%s*!", place.Name)))
	case errors.Is(err, lunch.ErrNoRerolls):
		return response.Ephemral(response.Section(response.PlainText("You don't have any more rerolls this week")))
	case errors.Is(err, lunch.ErrNoPlaces):
		return response.Ephemral(response.Section(response.PlainText("No places to choose from, add some!")))
	default:
		return response.InternalServerError(err)
	}
}

func handleAdd(ctx context.Context, place string) (*events.APIGatewayProxyResponse, error) {
	if err := roller.NewPlace(ctx, place); err != nil {
		return response.InternalServerError(err)
	}
	return response.Ephemral(response.Section(response.Markdown("*%s* added!", place)))
}
