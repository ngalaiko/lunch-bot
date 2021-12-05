package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

type ChallangeRequest struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func ParseRequest(r *http.Request) (*CommandRequest, *ActionsRequest, *ChallangeRequest, error) {
	ct := r.Header.Get("Content-Type")
	switch ct {
	case "application/x-www-form-urlencoded":
		r.ParseForm()
		if payload := r.Form.Get("payload"); payload != "" {
			actionsRequest := &ActionsRequest{}
			if err := json.NewDecoder(strings.NewReader(payload)).Decode(actionsRequest); err != nil {
				return nil, nil, nil, fmt.Errorf("failed to unmarshal payload as json: %w", err)
			}
			return nil, actionsRequest, nil, nil
		} else {
			return &CommandRequest{
				Command:  r.Form.Get("command"),
				Text:     r.Form.Get("text"),
				UserID:   r.Form.Get("user_id"),
				UserName: r.Form.Get("user_name"),
			}, nil, nil, nil
		}
	case "application/json":
		challangeReq := &ChallangeRequest{}
		if err := json.NewDecoder(r.Body).Decode(challangeReq); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal payload as json: %w", err)
		}
		return nil, nil, challangeReq, nil
	default:
		return nil, nil, nil, fmt.Errorf("unsupported content type: %s", ct)
	}
}
