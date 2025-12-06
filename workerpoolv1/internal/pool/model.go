package pool

import "time"

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID   int
	Outcome string
	Cost    time.Duration
	Err     string
}
