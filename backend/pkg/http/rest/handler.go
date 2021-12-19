package rest

import (
	"net/http"

	"lunch/pkg/http/rest/users"

	"github.com/go-chi/chi/v5"
)

func Handler() http.HandlerFunc {
	r := chi.NewRouter()
	r.Mount("/users", users.Handler())
	return r.ServeHTTP
}
