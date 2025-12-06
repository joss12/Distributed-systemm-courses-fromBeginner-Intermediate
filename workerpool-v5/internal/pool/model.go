package pool

import "time"

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID   int
	Outcome string
	Cost    time.Duration //simulated or measure time
	Err     string
}
