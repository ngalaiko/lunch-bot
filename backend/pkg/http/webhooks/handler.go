package webhooks

import (
	"fmt"
	"net/http"

	"lunch/pkg/http/webhooks/slack"
	"lunch/pkg/lunch"
	service_users "lunch/pkg/users/service"

	"github.com/go-chi/chi/v5"
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

func Handler(cfg *Configuration, roller *lunch.Roller, usersService *service_users.Service) http.Handler {
	r := chi.NewMux()
	r.Mount("/slack", slack.NewHandler(cfg.Slack, roller, usersService))
	return r
}
