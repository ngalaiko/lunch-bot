package users

import (
	"encoding/json"
	"log"
	"net/http"

	"lunch/pkg/users"

	"github.com/go-chi/chi/v5"
)

func getMe() http.HandlerFunc {
	type response struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := users.FromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("[ERROR] failed to encode response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Handler() http.HandlerFunc {
	r := chi.NewRouter()
	r.Get("/me", getMe())
	return r.ServeHTTP
}
