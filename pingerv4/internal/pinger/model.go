package pinger

import "time"

type Result struct {
	URL       string        `json:"url"`
	Status    int           `json:"status"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error"`
	Timestamp time.Time     `json:"timestamp"`
}

// Config groups runtime knows so we validate & pass as one unit
type Config struct {
	Timeout     time.Duration
	Concurrency int
}
