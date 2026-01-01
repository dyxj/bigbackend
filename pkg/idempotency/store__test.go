package idempotency

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLockConfig(t *testing.T) {
	config := DefaultLockConfig()

	assert.Equal(t, 20*time.Second, config.Expiry)
	assert.False(t, config.ShouldRetry)
	assert.Equal(t, 0, config.RetryAttempts)
	assert.Equal(t, time.Duration(0), config.RetryDelay)
}

func TestWithLockExpiry(t *testing.T) {
	config := &LockConfig{}

	WithLockExpiry(30 * time.Second)(config)

	assert.Equal(t, 30*time.Second, config.Expiry)
}

func TestWithLockRetry(t *testing.T) {
	config := &LockConfig{}

	WithLockRetry(5, 2*time.Second)(config)

	assert.True(t, config.ShouldRetry)
	assert.Equal(t, 5, config.RetryAttempts)
	assert.Equal(t, 2*time.Second, config.RetryDelay)
}

func TestWithLockNoRetry(t *testing.T) {
	config := &LockConfig{
		ShouldRetry:   true,
		RetryAttempts: 3,
		RetryDelay:    1 * time.Second,
	}

	WithLockNoRetry()(config)

	assert.False(t, config.ShouldRetry)
	assert.Equal(t, 0, config.RetryAttempts)
	assert.Equal(t, time.Duration(0), config.RetryDelay)
}
