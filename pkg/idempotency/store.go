package idempotency

import (
	"sync"
)

type Store interface {
	// Get returns the stored response for the idempotency key, or ok=false if not found.
	Get(key string) (resp *Response, ok bool)
	// Set stores the response for the idempotency key.
	Set(key string, resp *Response)
}

type MemStore struct {
	mu   sync.RWMutex
	data map[string]*Response
}

func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]*Response)}
}

func (s *MemStore) Get(key string) (*Response, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resp, ok := s.data[key]
	return resp, ok
}

func (s *MemStore) Set(key string, resp *Response) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = resp
}
