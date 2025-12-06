package main

import (
	"log"
	"net/http"
	"time"

	httpHandler "github.com/kvstore-v2/internal/http"
	"github.com/kvstore-v2/internal/middleware"
	"github.com/kvstore-v2/internal/store"
)

func main() {
	s := store.NewStore()
	h := httpHandler.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /kv/", h.Get)
	mux.HandleFunc("PUT /kv/", h.Put)
	mux.HandleFunc("DELETE /kv/", h.Delete)

	server := &http.Server{
		Addr:         ":8087",
		Handler:      middleware.Timeout(2*time.Second, mux),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Println("KV Store running on :8087")
	log.Fatal(server.ListenAndServe())
}
