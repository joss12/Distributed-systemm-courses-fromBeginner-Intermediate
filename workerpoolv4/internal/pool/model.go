package pool

import "time"

// Job is a task the worker must process
// We keep it generic: just an ID and any paypload.
type Job struct {
	ID      int
	Payload string
}

// Result is produced by workers.
type Result struct {
	JobID   int
	Outcome string
	Cost    time.Duration //simulate or measured time
	Err     string
}
