package http

import (
	"net/http"

	"lunch/pkg/http/oauth"
	"lunch/pkg/http/rest"
	"lunch/pkg/http/slack"
	"lunch/pkg/http/websocket"
	"lunch/pkg/lunch"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewHandler(cfg *Configuration, roller *lunch.Roller) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			// development
			"https://localhost:3000",
			"https://localhost:3001",
		},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.With(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded")).
		Post("/slack-lunch-bot", slack.NewHandler(cfg.Slack, roller).ServeHTTP)
	r.Get("/ws", websocket.Handler(roller).ServeHTTP)
	r.Mount("/oauth", oauth.Handler(cfg.OAuth))
	r.Mount("/api", rest.Handler())

	return r
}
