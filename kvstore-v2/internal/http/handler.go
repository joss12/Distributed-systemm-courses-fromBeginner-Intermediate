package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kvstore-v2/internal/response"
	"github.com/kvstore-v2/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func extractkey(path string) string {
	return strings.TrimPrefix(path, "/kv/")
}

// PUT /kv/{key}
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	key := extractkey(r.URL.Path)

	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, 400, "invalid json")
	}

	h.store.Put(key, body.Value)
	response.JSON(w, 200, map[string]string{"status": "ok"})
}

// GET /kv/{key}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := extractkey(r.URL.Path)

	v, err := h.store.Get(key)
	if err != nil {
		response.Error(w, 404, "key not found")
		return
	}
	response.JSON(w, 400, map[string]string{"value": v})
}

// DELETE /kv/{key}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := extractkey(r.URL.Path)

	if err := h.store.Delete(key); err != nil {
		response.Error(w, 404, "key not found")
		return
	}
	w.WriteHeader(204)
}
