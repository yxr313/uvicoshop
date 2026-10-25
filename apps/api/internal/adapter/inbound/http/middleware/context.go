package middleware

import "context"

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUserID).(string)
	return v, ok && v != ""
}

func MustUserID(ctx context.Context) (string, error) {
	v, ok := UserIDFromContext(ctx)
	if !ok {
		return "", errMissingUser
	}
	return v, nil
}