package auth

import (
	"net/http"
	"time"

	"lunch/pkg/jwt"
)

const (
	cookieName = "auth"
)

func RemoveCookie(w http.ResponseWriter, secure bool) {
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "",
	}

	http.SetCookie(w, &cookie)
}

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
