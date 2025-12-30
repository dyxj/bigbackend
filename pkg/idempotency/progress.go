package idempotency

import (
	"context"
	"errors"
)

var ErrInProgress = errors.New("idempotency key in progress")

type contextKey string

const isInProgressKey contextKey = "is-idempotency-in-progress"

func SetInProgress(ctx context.Context, v bool) {
	context.WithValue(ctx, isInProgressKey, v)
}
