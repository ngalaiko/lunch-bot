package auth

import (
	"net/http"
	"time"

	"lunch/pkg/jwt"
)

const (
	cookieName = "auth"
)

func SetCookie(w http.ResponseWriter, token *jwt.Token, secure bool) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    token.Token,
		MaxAge:   int(time.Until(token.ExpiresAt).Seconds()),
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "",
	}

	http.SetCookie(w, &cookie)
}
