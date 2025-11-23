package pinger

import "time"

// Result is a singe piing outcome. Keep it small and serailizable.
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
