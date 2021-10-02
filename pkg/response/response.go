package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type response struct {
	ResponseType responseType `json:"response_type"`
	Text         string       `json:"text"`
}

type responseType string

const (
	responseTypeEphemeral responseType = "ephemeral"
	responseTypeInChannel responseType = "in_channel"
)

// Ephemral sends a message back visible only by the caller.
func Ephemral(text string) (*events.APIGatewayProxyResponse, error) {
	return respondJSON(&response{
		ResponseType: responseTypeEphemeral,
		Text:         text,
	})
}

// InChannel sends a message back visible by everyone in the channel.
func InChannel(text string) (*events.APIGatewayProxyResponse, error) {
	return respondJSON(&response{
		ResponseType: responseTypeInChannel,
		Text:         text,
	})
}

func BadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	return Ephemral(err.Error())
}

func InternalServerError(err error) (*events.APIGatewayProxyResponse, error) {
	log.Printf("[ERROR] %s", err)
	return Ephemral("Sorry, that didn't work. Try again or contact the app administrator.")
}

func respondJSON(body interface{}) (*events.APIGatewayProxyResponse, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return InternalServerError(fmt.Errorf("failed to marshal response: %w", err))
	}
	return &events.APIGatewayProxyResponse{
		Body: string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
	}, nil
}
