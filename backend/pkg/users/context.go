package users

import "context"

type userContextKey struct{}

var contextKey = userContextKey{}

func NewContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextKey, user)
}

func FromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(contextKey).(*User)
	return user, ok
}
