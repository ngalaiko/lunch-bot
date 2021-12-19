package users

import (
	"encoding/json"
	"log"
	"net/http"

	"lunch/pkg/jwt"

	"github.com/go-chi/chi/v5"
)

func getMe() http.HandlerFunc {
	type response struct {
		ID string `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := jwt.FromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp := &response{
			ID: token.UserID,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
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
