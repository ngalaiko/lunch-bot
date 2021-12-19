package oauth

import (
	"fmt"
	"net/http"

	"lunch/pkg/http/oauth/slack"
	"lunch/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func Handler(cfg *Configuration, jwtService *jwt.Service) http.HandlerFunc {
	r := chi.NewRouter()
	applicationJSON := middleware.AllowContentType("application/json")
	r.With(applicationJSON).Post("/slack", slack.Handler(cfg.Slack, jwtService))
	return r.ServeHTTP
}
