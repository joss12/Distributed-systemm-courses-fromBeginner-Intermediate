package main

import (
	"log"
	"os"
	"strings"

	"github.com/kv-repl/internal/kv"
)

func main() {
	addr := os.Getenv("ADDR")
	peers := strings.Split(os.Getenv("PEERS"), ",")

	if addr == "" {
		addr = "127.0.0.1:9001"
	}
	if len(peers) == 1 && peers[0] == "" {
		peers = []string{addr}
	}

	node := kv.NewNode(addr, peers)
	log.Fatal(node.Run())
}
