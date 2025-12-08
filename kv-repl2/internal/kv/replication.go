package kv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Replicator struct {
	Peers  []string
	Self   string
	Client *http.Client
}

type ReplMsg struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *Replicator) Broadcast(msg ReplMsg) {
	fmt.Println("[Broadcast] firing:", msg)

	body, _ := json.Marshal(msg)

	for _, peer := range r.Peers {
		if peer == r.Self {
			continue
		}

		go func(p string) {
			fmt.Println("[Broadcast] sending to peer:", p)

			//Replication context must be indepenedent of the HTTP request
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			req, _ := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				"http://"+p+"/replicator",
				bytes.NewReader(body),
			)
			req.Header.Set("Content-Type", "application/json")

			resp, err := r.Client.Do(req)
			if err != nil {
				fmt.Println("[Broadcast] ERROR sending to", p, ":", err)
				return
			}
			resp.Body.Close()

			fmt.Println("[Broadcast] SUCCESS sending to", p)
		}(peer)
	}
}
