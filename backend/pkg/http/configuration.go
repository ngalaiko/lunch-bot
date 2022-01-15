package http

import (
	"fmt"

	"lunch/pkg/http/oauth"
	"lunch/pkg/http/webhooks"
)

type Configuration struct {
	Webhooks *webhooks.Configuration
	OAuth    *oauth.Configuration
}

func (c *Configuration) Parse() error {
	c.Webhooks = &webhooks.Configuration{}
	if err := c.Webhooks.Parse(); err != nil {
		return fmt.Errorf("failed to parse slack configuration: %w", err)
	}

	c.OAuth = &oauth.Configuration{}
	if err := c.OAuth.Parse(); err != nil {
		return fmt.Errorf("failed to parse oauth configuration: %w", err)
	}

	return nil
}
