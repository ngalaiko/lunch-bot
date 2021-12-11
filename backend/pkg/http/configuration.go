package http

import (
	"fmt"

	"lunch/pkg/http/slack"
)

type Configuration struct {
	Slack *slack.Configuration
}

func (c *Configuration) Parse() error {
	c.Slack = &slack.Configuration{}
	if err := c.Slack.Parse(); err != nil {
		return fmt.Errorf("failed to parse slack configuration: %w", err)
	}
	return nil
}
