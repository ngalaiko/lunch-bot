package jwt

import (
	"time"

	"lunch/pkg/users"
)

// Token contains information about an authorized user.
type Token struct {
	Token     string      `json:"token"`
	User      *users.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}
