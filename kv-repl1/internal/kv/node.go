package kv

import (
	"log"
	"net/http"
	"time"
)

type Node struct {
	Addr  string
	Peers []string
	Store *Store
	Repl  *Replicator
	Http  *http.Server
}

func NewNode(addr string, peers []string) *Node {
	store := NewStore()
	repl := &Replicator{
		Peers:  peers,
		Self:   addr,
		Client: &http.Client{Timeout: 2 * time.Second},
	}

	h := &NodeHandlers{Store: store, Replicator: repl}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /kv/", h.HandlePut)
	mux.HandleFunc("GET /kv/", h.HandleGet)
	mux.HandleFunc("GET /replicate", h.HandleReplicate)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	mux.HandleFunc("/kv/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.HandlePut(w, r)
			return
		}
		if r.Method == http.MethodGet {
			h.HandleGet(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/replicate", h.HandleReplicate)

	s := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	return &Node{
		Addr:  addr,
		Peers: peers,
		Store: store,
		Repl:  repl,
		Http:  s,
	}
}

func (n *Node) Run() error {
	log.Println("node listening on", n.Addr, "peers:", n.Peers)
	return n.Http.ListenAndServe()
}
