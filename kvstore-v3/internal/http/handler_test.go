package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpHandler "github.com/kvstore-v3/internal/http"
	"github.com/kvstore-v3/internal/store"
)

func newTestServer() (*http.ServeMux, *store.Store) {
	s := store.NewStore()
	h := httpHandler.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /kv/", h.Get)
	mux.HandleFunc("PUT /kv/", h.Put)
	mux.HandleFunc("DELETE /kv/", h.Delete)

	return mux, s
}

func TestPutHandler(t *testing.T) {
	mux, _ := newTestServer()

	body := bytes.NewBufferString(`{"value":"golang"}`)
	req := httptest.NewRequest("PUT", "/kv/lang", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("expected status 200, go %d", rr.Code)
	}
}

func TestGetHandler(t *testing.T) {
	mux, s := newTestServer()
	s.Put("framework", "fiber")

	req := httptest.NewRequest("GET", "/kv/framework", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["value"] != "fiber" {
		t.Fatalf("invalid value='fiber', got '%s'", resp["value"])
	}
}

func TestDeleteHandler(t *testing.T) {
	mux, s := newTestServer()
	s.Put("name", "eddy")

	req := httptest.NewRequest("DELETE", "/kv/name", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != 204 {
		t.Fatalf("expected status 204, go %d", rr.Code)
	}

	//ensure deleted
	_, err := s.Get("name")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
