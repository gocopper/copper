package cauth

import "context"

type ctxKey string

const (
	userCtxKey = ctxKey("cauth/user")
)

// GetCurrentUser returns the current authenticated user from the context.
// Returns nil if there is no user in the context.
func GetCurrentUser(ctx context.Context) *user {
	user, ok := ctx.Value(userCtxKey).(*user)
	if !ok {
		return nil
	}

	return user
}

func ctxWithUser(ctx context.Context, user *user) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}
