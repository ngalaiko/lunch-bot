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
