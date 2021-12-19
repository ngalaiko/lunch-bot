package http

import (
	"fmt"

	"lunch/pkg/http/oauth"
	"lunch/pkg/http/slack"
)

type Configuration struct {
	Slack *slack.Configuration
	OAuth *oauth.Configuration
}

func (c *Configuration) Parse() error {
	c.Slack = &slack.Configuration{}
	if err := c.Slack.Parse(); err != nil {
		return fmt.Errorf("failed to parse slack configuration: %w", err)
	}

	c.OAuth = &oauth.Configuration{}
	if err := c.OAuth.Parse(); err != nil {
		return fmt.Errorf("failed to parse oauth configuration: %w", err)
	}

	return nil
}
