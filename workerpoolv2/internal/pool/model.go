package pool

import "time"

type Job struct {
	ID      int
	Payload string
}

// Result is produced by worker.
type Result struct {
	JobID   int
	Outcome string
	Cost    time.Duration //simulated or measured time
	Err     string
}
