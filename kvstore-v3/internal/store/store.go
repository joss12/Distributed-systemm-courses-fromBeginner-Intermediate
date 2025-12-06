package store

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Store struct {
	mu sync.RWMutex
	m  map[string]string
}

// NewStore initiaalizes the map safely.
func NewStore() *Store {
	return &Store{
		m: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.m[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *Store) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.m[key]; !exists {
		return ErrNotFound
	}
	delete(s.m, key)
	return nil
}
