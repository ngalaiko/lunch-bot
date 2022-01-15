package http

import (
	"net/http"

	"lunch/pkg/http/auth"
	"lunch/pkg/http/oauth"
	"lunch/pkg/http/rest"
	"lunch/pkg/http/webhooks"
	"lunch/pkg/http/websocket"
	"lunch/pkg/jwt"
	"lunch/pkg/lunch"
	"lunch/pkg/lunch/events"
	service_users "lunch/pkg/users/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewHandler(
	cfg *Configuration,
	roller *lunch.Roller,
	jwtService *jwt.Service,
	usersService *service_users.Service,
	eventsRegistry *events.Registry,
) http.Handler {
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
	r.Use(auth.Parser(jwtService))

	r.Route("/api", func(r chi.Router) {
		r.Mount("/webhooks", webhooks.Handler(cfg.Webhooks, roller, usersService))
		r.Mount("/oauth", oauth.Handler(cfg.OAuth, jwtService, usersService))
		r.Mount("/ws", websocket.Handler(roller, eventsRegistry))
		r.Mount("/", rest.Handler())
	})

	return r
}
