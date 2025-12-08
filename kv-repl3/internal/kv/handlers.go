package kv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type NodeHandlers struct {
	Store      *Store
	Replicator *Replicator
}

type putReq struct {
	Value string `json:"value"`
}

type putRes struct {
	Status string `json:"status"`
}

func (h *NodeHandlers) HandlePut(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")

	var req putReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "had json", 400)
		return
	}

	id := generateID()
	applied := h.Store.Apply(id, key, req.Value)

	if applied {
		h.Replicator.Broadcast(ReplMsg{
			ID:    id,
			Key:   key,
			Value: req.Value,
		})
	}
	_ = json.NewEncoder(w).Encode(putRes{Status: "ok"})
}

func (h *NodeHandlers) HandleGet(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")
	if v, ok := h.Store.Get(key); ok {
		_ = json.NewEncoder(w).Encode(map[string]string{"value": v})
		return
	}
	http.NotFound(w, r)
}
func (h *NodeHandlers) HandleReplicate(w http.ResponseWriter, r *http.Request) {
	var msg ReplMsg
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		fmt.Println("[Replicate] bad json:", err)
		http.Error(w, "bad json", 400)
		return
	}

	fmt.Println("[Replicate] recieved msg:", msg)
	h.Store.Apply(msg.ID, msg.Key, msg.Value)
	w.WriteHeader(200)
}

func generateID() string {
	return uuid.NewString()
}
