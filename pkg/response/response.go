package response

import (
	"bytes"
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

// EphemralURL sends a message back via response URL.
func EphemralURL(responseURL string, sections ...*Block) (*events.APIGatewayProxyResponse, error) {
	return respondJSONURL(responseURL, newMessage(responseTypeEphemeral, sections...))
}

// ReplaceEphemralURL sends a message back via response URL, replacing original message.
func ReplaceEphemralURL(responseURL string, sections ...*Block) (*events.APIGatewayProxyResponse, error) {
	return respondJSONURL(responseURL, newReplaceMessage(responseTypeEphemeral, sections...))
}

// Ephemral sends a message back visible only by the caller.
func Ephemral(sections ...*Block) (*events.APIGatewayProxyResponse, error) {
	return respondJSON(newMessage(responseTypeEphemeral, sections...))
}

// InChannel sends a message back visible by everyone in the channel.
func InChannel(sections ...*Block) (*events.APIGatewayProxyResponse, error) {
	return respondJSON(newMessage(responseTypeInChannel, sections...))
}

func BadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	return Ephemral(Section(PlainText(err.Error())))
}

func InternalServerError(err error) (*events.APIGatewayProxyResponse, error) {
	log.Printf("[ERROR] %s", err)
	return Ephemral(Section(PlainText("Sorry, that didn't work. Try again or contact the app administrator.")))
}

func respondJSONURL(responseURL string, body interface{}) (*events.APIGatewayProxyResponse, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return InternalServerError(fmt.Errorf("failed to marshal response: %w", err))
	}
	log.Printf("[TRACE] response: %s", string(data))
	resp, err := http.Post(responseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return InternalServerError(fmt.Errorf("failed to post message to slack: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		return InternalServerError(fmt.Errorf("got invalid status from slack: %d", resp.StatusCode))
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func respondJSON(body interface{}) (*events.APIGatewayProxyResponse, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return InternalServerError(fmt.Errorf("failed to marshal response: %w", err))
	}
	log.Printf("[TRACE] response: %s", string(bytes))
	return &events.APIGatewayProxyResponse{
		Body: string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
	}, nil
}
