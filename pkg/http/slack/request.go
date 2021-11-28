package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hetiansu5/urlquery"
)

var signingSecret = "notset"

type FormURLEncodedRequest struct {
	Command  string `query:"command"`
	Text     string `query:"text"`
	UserID   string `query:"user_id"`
	UserName string `query:"user_name"`

	Payload string `query:"payload"`
}

type CommandRequest struct {
	Command  string
	Text     string
	UserID   string
	UserName string
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Action struct {
	ActionID string `json:"action_id"`
	Value    string `json:"value"`
}

type ActionsRequest struct {
	User        *User     `json:"user"`
	Actions     []*Action `json:"actions"`
	ResponseUrl string    `json:"response_url"`
}

func ParseRequest(r *http.Request) (*CommandRequest, *ActionsRequest, error) {
	ct := r.Header.Get("Content-Type")
	if ct != "application/x-www-form-urlencoded" {
		return nil, nil, fmt.Errorf("invalid content type: '%s', expected 'application/x-www-form-urlencoded'", ct)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read request body: %w", err)
	}

	req := &FormURLEncodedRequest{}
	if err := urlquery.Unmarshal(body, req); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal as query: %w", err)
	}

	if req.Payload != "" {
		actionsRequest := &ActionsRequest{}
		if err := json.NewDecoder(strings.NewReader(req.Payload)).Decode(actionsRequest); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal as json: %w", err)
		}
		return nil, actionsRequest, nil
	}

	return &CommandRequest{
		Command:  req.Command,
		Text:     req.Text,
		UserID:   req.UserID,
		UserName: req.UserName,
	}, nil, nil
}
