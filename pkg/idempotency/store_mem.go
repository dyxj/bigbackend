package idempotency

import (
	"context"
	"sync"
	"time"
)

type MemStore struct {
	muData sync.RWMutex
	data   map[string]*Response

	muLocks         sync.RWMutex
	locks           map[string]struct{}
	defLockConfigFn func() *LockConfig
}

func NewMemStore(defLockConfigFn func() *LockConfig) *MemStore {
	return &MemStore{
		data:            make(map[string]*Response),
		locks:           make(map[string]struct{}),
		defLockConfigFn: defLockConfigFn,
	}
}

func (s *MemStore) Lock(ctx context.Context, key string, opts ...LockOptions) error {
	config := s.defLockConfigFn()

	for _, opt := range opts {
		opt(config)
	}

	retries := 0
	for {
		_, isLocked := s.getLock(key)
		if !isLocked {
			s.lock(ctx, key)
			// set lock expiry
			if config.Expiry > 0 {
				time.AfterFunc(config.Expiry, func() {
					_ = s.Unlock(context.Background(), key)
				})
			}
			return nil // lock obtained
		}

		if !config.ShouldRetry || retries > config.RetryAttempts {
			return ErrInProgress
		}

		retries++

		select {
		case <-ctx.Done():
			return ErrInProgress
		case <-time.After(config.RetryDelay):
		}
	}
}

func (s *MemStore) getLock(key string) (struct{}, bool) {
	s.muLocks.RLock()
	defer s.muLocks.RUnlock()
	sKey, ok := s.locks[key]
	return sKey, ok
}

func (s *MemStore) lock(ctx context.Context, key string) {
	s.muLocks.Lock()
	defer s.muLocks.Unlock()
	s.locks[key] = struct{}{}
}

func (s *MemStore) Unlock(ctx context.Context, key string) error {
	s.muLocks.Lock()
	defer s.muLocks.Unlock()
	delete(s.locks, key)
	return nil
}

func (s *MemStore) Get(ctx context.Context, key string) (*Response, error) {
	s.muData.RLock()
	defer s.muData.RUnlock()
	resp, ok := s.data[key]
	if !ok {
		return nil, nil
	}
	return resp, nil
}

func (s *MemStore) Set(ctx context.Context, key string, resp *Response, expiry time.Duration) error {
	s.muData.Lock()
	defer s.muData.Unlock()
	s.data[key] = resp
	if expiry > 0 {
		time.AfterFunc(expiry, func() {
			s.unset(key)
		})
	}
	return nil
}

func (s *MemStore) unset(key string) {
	s.muData.Lock()
	defer s.muData.Unlock()
	delete(s.data, key)
}
