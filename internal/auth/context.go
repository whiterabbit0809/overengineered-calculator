// internal/auth/context.go
package auth

import "context"

// contextKey is an unexported type to avoid key collisions.
type contextKey string

const (
	ctxKeyUserID contextKey = "userID"
	ctxKeyEmail  contextKey = "email"
)

// ContextWithUser stores userID and email in the context.
// Call this in middleware after verifying the JWT.
func ContextWithUser(ctx context.Context, userID, email string) context.Context {
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyEmail, email)
	return ctx
}

// UserFromContext extracts userID and email from the context.
// Call this in handlers that need to know who the current user is.
func UserFromContext(ctx context.Context) (userID, email string, ok bool) {
	uid, ok1 := ctx.Value(ctxKeyUserID).(string)
	em, ok2 := ctx.Value(ctxKeyEmail).(string)
	if !ok1 || !ok2 {
		return "", "", false
	}
	return uid, em, true
}
