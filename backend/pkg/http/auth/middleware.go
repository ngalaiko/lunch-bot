package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"lunch/pkg/jwt"
	"lunch/pkg/users"
)

const authLeeway = time.Hour * 24

func Parser(jwtService *jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := fromCookie(r.Cookies())
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			jwt, err := jwtService.Verify(r.Context(), token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if time.Until(jwt.ExpiresAt) < authLeeway {
				jwt, err = jwtService.NewToken(r.Context(), jwt.User)
				if err != nil {
					log.Printf("[ERROR] failed to generate new token: %s", err)
					next.ServeHTTP(w, r)
					return
				}
				secure := r.TLS != nil
				SetCookie(w, jwt, secure)
			}

			ctx := users.NewContext(r.Context(), jwt.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func fromCookie(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("no cookie found")
}
