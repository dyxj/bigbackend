package idempotency

import (
	"context"
	"time"
)

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

func WithLockExpiry(expiry time.Duration) LockOptions {
	return func(config *LockConfig) {
		config.Expiry = expiry
	}
}

func WithLockRetry(attempts int, delay time.Duration) LockOptions {
	return func(config *LockConfig) {
		config.ShouldRetry = true
		config.RetryAttempts = attempts
		config.RetryDelay = delay
	}
}

func WithLockNoRetry() LockOptions {
	return func(config *LockConfig) {
		config.ShouldRetry = false
		config.RetryAttempts = 0
		config.RetryDelay = 0
	}
}
