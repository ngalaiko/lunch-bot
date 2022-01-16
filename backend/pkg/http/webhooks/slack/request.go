package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func verifySlackSignatureV0(header http.Header, body []byte, signingSecret string) error {
	signature := header.Get("X-Slack-Signature")
	if signature == "" {
		return fmt.Errorf("X-Slack-Signature is missing")
	}

	mac := header.Get("X-Slack-Request-Timestamp")
	if mac == "" {
		return fmt.Errorf("X-Slack-Request-Timestamp is missing")
	}

	baseString := fmt.Sprintf("v0:%s:%s", mac, string(body))
	hash := hmac.New(sha256.New, []byte(signingSecret))
	if _, err := hash.Write([]byte(baseString)); err != nil {
		return fmt.Errorf("failed to calculate signature: %w", err)
	}

	expected := fmt.Sprintf("v0=%x", hash.Sum(nil))
	if expected != signature {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func ParseRequest(r *http.Request, signingSecret string) (*CommandRequest, *ActionsRequest, *ChallangeRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read request body: %w", err)
	}

	if len(signingSecret) > 0 {
		if err := verifySlackSignatureV0(r.Header, body, signingSecret); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to verify request signature: %w", err)
		}
	}

	ct := r.Header.Get("Content-Type")
	switch ct {
	case "application/x-www-form-urlencoded":
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse request body: %w", err)
		}
		if payload := values.Get("payload"); payload != "" {
			actionsRequest := &ActionsRequest{}
			if err := json.NewDecoder(strings.NewReader(payload)).Decode(actionsRequest); err != nil {
				return nil, nil, nil, fmt.Errorf("failed to unmarshal payload as json: %w", err)
			}
			return nil, actionsRequest, nil, nil
		} else {
			return &CommandRequest{
				Command:  values.Get("command"),
				Text:     values.Get("text"),
				UserID:   values.Get("user_id"),
				UserName: values.Get("user_name"),
			}, nil, nil, nil
		}
	case "application/json":
		challangeReq := &ChallangeRequest{}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(challangeReq); err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal payload as json: %w", err)
		}
		return nil, nil, challangeReq, nil
	default:
		return nil, nil, nil, fmt.Errorf("unsupported content type: %s", ct)
	}
}
