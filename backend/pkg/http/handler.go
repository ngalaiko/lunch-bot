package http

import (
	"net/http"

	"lunch/pkg/http/slack"
	"lunch/pkg/http/websocket"
	"lunch/pkg/lunch"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(cfg *Configuration, roller *lunch.Roller) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)

	r.With(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded")).
		Post("/slack-lunch-bot", slack.NewHandler(cfg.Slack, roller).ServeHTTP)
	r.Get("/ws", websocket.Handler(roller).ServeHTTP)

	return r
}
