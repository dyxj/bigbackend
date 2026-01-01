package idempotency

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemStore_Lock_LockObtained(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	err := store.Lock(context.Background(), "test-key")
	assert.NoError(t, err)
	err2 := store.Lock(context.Background(), "test-key2")
	assert.NoError(t, err2)

	assert.Contains(t, store.locks, "test-key")
	assert.Contains(t, store.locks, "test-key2")
}

func TestMemStore_Lock_LockExpired(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	err := store.Lock(context.Background(), "test-key",
		WithLockExpiry(100*time.Millisecond),
	)
	assert.NoError(t, err)
	assert.Contains(t, store.locks, "test-key")

	<-time.After(150 * time.Millisecond)

	assert.NotContains(t, store.locks, "test-key")
}

func TestMemStore_Lock_ErrInProgress(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	err := store.Lock(context.Background(), "test-key")
	assert.NoError(t, err)
	err2 := store.Lock(context.Background(), "test-key")
	assert.ErrorIs(t, err2, ErrInProgress)
}

func TestMemStore_Lock_WithRetry_ErrInProgress(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{
			Expiry:        0,
			ShouldRetry:   false,
			RetryAttempts: 0,
			RetryDelay:    0,
		}
	})

	err := store.Lock(context.Background(), "test-key")
	assert.NoError(t, err)

	start := time.Now()
	err2 := store.Lock(context.Background(), "test-key",
		WithLockRetry(3, 100*time.Millisecond),
	)
	duration := time.Since(start)

	assert.ErrorIs(t, err2, ErrInProgress)
	assert.GreaterOrEqual(t, duration, 300*time.Millisecond)
}

func TestMemStore_Lock_WithRetry_ErrInProgress_ContextDone(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{
			Expiry:        0,
			ShouldRetry:   false,
			RetryAttempts: 0,
			RetryDelay:    0,
		}
	})

	err := store.Lock(context.Background(), "test-key")
	assert.NoError(t, err)

	ctx, cancelFn := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancelFn()
	start := time.Now()
	err2 := store.Lock(ctx, "test-key",
		WithLockRetry(5, 100*time.Millisecond),
	)
	duration := time.Since(start)

	assert.ErrorIs(t, err2, ErrInProgress)
	assert.GreaterOrEqual(t, duration, 100*time.Millisecond)
	assert.Less(t, duration, 300*time.Millisecond)
}

func TestMemStore_Lock_WithRetry_LockObtained(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{
			Expiry:        0,
			ShouldRetry:   false,
			RetryAttempts: 0,
			RetryDelay:    0,
		}
	})

	err := store.Lock(context.Background(), "test-key",
		WithLockExpiry(200*time.Millisecond),
	)
	assert.NoError(t, err)

	start := time.Now()
	err2 := store.Lock(context.Background(), "test-key",
		WithLockRetry(5, 100*time.Millisecond),
	)
	duration := time.Since(start)

	assert.NoError(t, err2)
	assert.GreaterOrEqual(t, duration, 200*time.Millisecond)
}

func TestMemStore_Unlock(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	err := store.Lock(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.Contains(t, store.locks, "test-key")

	err = store.Unlock(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.NotContains(t, store.locks, "test-key")
}

func TestMemStore_Set_Get(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	resp := &Response{}

	err := store.Set(context.Background(), "test-key", resp, 0)
	assert.NoError(t, err)

	result, err := store.Get(context.Background(), "test-key")
	assert.NoError(t, err)

	assert.Equal(t, resp, result)
}

func TestMemStore_Set_Expired(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	resp := &Response{}

	err := store.Set(context.Background(), "test-key", resp, 100*time.Millisecond)
	assert.NoError(t, err)

	<-time.After(150 * time.Millisecond)

	result, err := store.Get(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestMemStore_Get_NotFound(t *testing.T) {
	t.Parallel()
	store := NewMemStore(func() *LockConfig {
		return &LockConfig{}
	})

	result, err := store.Get(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.Nil(t, result)
}
