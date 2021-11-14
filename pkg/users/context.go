package users

import "context"

type userContextKey struct{}

var contextKey = userContextKey{}

type User struct {
	Name string `dynamodbav:"name" json:"name"`
	ID   string `dynamodbav:"id" json:"id"`
}

func NewContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextKey, user)
}

func FromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(contextKey).(*User)
	return user, ok
}
