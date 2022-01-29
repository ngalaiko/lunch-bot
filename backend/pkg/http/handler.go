package http

import (
	"log"
	"net/http"
	"time"

	"lunch/pkg/http/auth"
	"lunch/pkg/http/oauth"
	"lunch/pkg/http/rest"
	"lunch/pkg/http/webhooks"
	"lunch/pkg/http/websocket"
	"lunch/pkg/jwt"
	"lunch/pkg/lunch"
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
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(&logFormatter{}))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/api/ping"))
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
		r.Mount("/ws", websocket.Handler(roller))
		r.Mount("/", rest.Handler())
	})

	return r
}

type logEntry struct {
	Path   string
	Method string
}

func (le *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	log.Printf("[INFO] %s %s %d %s", le.Method, le.Path, status, elapsed)
}

func (le *logEntry) Panic(v interface{}, stack []byte) {
	log.Printf("[ERROR] panic in handler: %+v\n%s", v, string(stack))
}

type logFormatter struct{}

func (lf *logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &logEntry{
		Path:   r.URL.Path,
		Method: r.Method,
	}
}
