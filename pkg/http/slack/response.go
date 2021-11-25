package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// EphemeralURL sends a message back via response URL.
func EphemeralURL(w http.ResponseWriter, responseURL string, text Text, sections ...*Block) error {
	return respondJSONURL(w, responseURL, newMessage(responseTypeEphemeral, text, sections...))
}

// ReplaceEphemeralURL sends a message back via response URL, replacing original message.
func ReplaceEphemeralURL(w http.ResponseWriter, responseURL string, text Text, sections ...*Block) error {
	return respondJSONURL(w, responseURL, newReplaceMessage(responseTypeEphemeral, text, sections...))
}

// Ephemeral sends a message back visible only by the caller.
func Ephemeral(w http.ResponseWriter, text Text, sections ...*Block) error {
	return respondJSON(w, newMessage(responseTypeEphemeral, text, sections...))
}

// InChannel sends a message back visible by everyone in the channel.
func InChannel(w http.ResponseWriter, text Text, sections ...*Block) error {
	return respondJSON(w, newMessage(responseTypeInChannel, text, sections...))
}

func BadRequest(w http.ResponseWriter, err error) error {
	return Ephemeral(w, Text(err.Error()), Section(PlainText(err.Error())))
}

func InternalServerError(w http.ResponseWriter, err error) error {
	log.Printf("[ERROR] %s", err)
	msg := "Sorry, that didn't work. Try again or contact the app administrator."
	return Ephemeral(w, Text(msg), Section(PlainText(msg)))
}

func respondJSONURL(w http.ResponseWriter, responseURL string, body interface{}) error {
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(body); err != nil {
		return InternalServerError(w, fmt.Errorf("failed to marshal response: %w", err))
	}
	resp, err := http.Post(responseURL, "application/json", b)
	if err != nil {
		return InternalServerError(w, fmt.Errorf("failed to post message to slack: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		return InternalServerError(w, fmt.Errorf("got invalid status from slack: %d", resp.StatusCode))
	}
	return nil
}

func respondJSON(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return InternalServerError(w, fmt.Errorf("failed to marshal response: %w", err))
	}
	return nil
}
