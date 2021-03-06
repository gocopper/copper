package cauth

import "context"

const ctxKeySession = "auth/session"

func GetCurrentUserUUID(ctx context.Context) string {
	userUUID, ok := ctx.Value(ctxKeySession).(string)
	if !ok || userUUID == "" {
		panic("user uuid not found in context")
	}

	return userUUID
}
