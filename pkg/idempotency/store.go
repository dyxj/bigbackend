package idempotency

import (
	"context"
	"errors"
	"time"
)

var ErrInProgress = errors.New("key in progress")

type Store interface {
	Lock(ctx context.Context, key string, opts ...LockOptions) error
	Unlock(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (*Response, error)
	Set(ctx context.Context, key string, resp *Response, expiry time.Duration) error
}

type LockConfig struct {
	Expiry        time.Duration
	ShouldRetry   bool
	RetryAttempts int
	RetryDelay    time.Duration
}

func DefaultLockConfig() *LockConfig {
	return &LockConfig{
		Expiry:        20 * time.Second,
		ShouldRetry:   false,
		RetryAttempts: 0,
		RetryDelay:    0,
	}
}

type LockOptions func(*LockConfig)

func WithExpiry(expiry time.Duration) LockOptions {
	return func(config *LockConfig) {
		config.Expiry = expiry
	}
}

func WithRetry(attempts int, delay time.Duration) LockOptions {
	return func(config *LockConfig) {
		config.ShouldRetry = true
		config.RetryAttempts = attempts
		config.RetryDelay = delay
	}
}
