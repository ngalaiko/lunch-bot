package http

import (
	"net/http"

	"lunch/pkg/lunch"
	storage_boosts "lunch/pkg/lunch/boosts/storage"
	storage_places "lunch/pkg/lunch/places/storage"
	storage_rolls "lunch/pkg/lunch/rolls/storage"
)

type Server struct {
	router *http.ServeMux

	roller *lunch.Roller
}

func NewServer(
	boostsStore storage_boosts.Storage,
	placesStore storage_places.Storage,
	rollsStore storage_rolls.Storage,
) *Server {
	s := &Server{
		router: http.NewServeMux(),
		roller: lunch.New(placesStore, boostsStore, rollsStore),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	normalizePath(accessLogs(s.router.ServeHTTP)).ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleSlack())
}
