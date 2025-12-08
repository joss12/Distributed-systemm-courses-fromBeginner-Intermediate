package kv

import (
	"log"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
	seen map[string]struct{} //write-id deduplication
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
		seen: make(map[string]struct{}),
	}
}

func (s *Store) Apply(id, key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.seen[id]; ok {
		log.Println("[Store] duplicate ignore:", id)
		return true
	}
	log.Println("[Store] apply:", id, key, value)

	s.seen[id] = struct{}{}
	s.data[key] = value
	return true
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}
