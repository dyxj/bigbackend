package idempotency

import "context"

type contextKey struct{}

// WithValue sets the idempotency key into the context
func WithValue(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, contextKey{}, key)
}

// FromContext returns value from context, empty string if not set
func FromContext(ctx context.Context) string {
	key, ok := ctx.Value(contextKey{}).(string)
	if !ok {
		return ""
	}
	return key
}
