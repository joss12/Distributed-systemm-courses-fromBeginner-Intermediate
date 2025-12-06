package pool

import (
	"context"
	"math/rand"
	"time"
)

func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			//context cancellation ends
			return

		case j, ok := <-jobs:
			if !ok {
				//Channel closed -> no more jobs
				return
			}
			//simulate processing time
			start := time.Now()
			d := time.Duration(20+rand.Intn(80)) * time.Millisecond
			time.Sleep(d)

			//simulate random failure
			if rand.Intn(10) == 0 {
				results <- Result{
					JobID: j.ID,
					Cost:  time.Since(start),
					Err:   "simulate",
				}
				continue
			}
			results <- Result{
				JobID:   j.ID,
				Cost:    time.Since(start),
				Outcome: "ok",
			}
		}
	}
}
