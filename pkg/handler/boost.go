package handler

import (
	"context"

	"lunch/pkg/response"

	"github.com/aws/aws-lambda-go/events"
)

func handeBoost(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	names, err := roller.ListPlaces(ctx)
	if err != nil {
		return response.InternalServerError(err)
	}

	options := make([]*response.Option, 0, len(names))
	for _, name := range names {
		options = append(options, &response.Option{
			Text:  response.PlainText(string(name)),
			Value: string(name),
		})
	}

	return response.Ephemral(
		response.Select(
			response.PlainText("Choose a place to boost:"),
			response.Static(
				response.PlainText("pick a place..."),
				"boost",
				options...,
			),
		),
	)
}
