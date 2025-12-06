package store_test

import (
	"testing"

	"github.com/kvstore/internal/store"
)

func TestPutAndGet(t *testing.T) {
	s := store.NewStore()
	s.Put("name", "golang")

	v, err := s.Get("name")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if v != "golang" {
		t.Fatalf("expected 'golang', got %s", v)
	}
}

func TestGetNotFound(t *testing.T) {
	s := store.NewStore()

	_, err := s.Get("missing")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	s := store.NewStore()
	s.Put("age", "30")

	err := s.Delete("age")
	if err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}

	_, err = s.Get("age")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteMissingKey(t *testing.T) {
	s := store.NewStore()

	err := s.Delete("missing-key")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, go %v", err)
	}
}
