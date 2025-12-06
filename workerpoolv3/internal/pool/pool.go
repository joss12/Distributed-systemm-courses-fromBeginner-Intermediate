package pool

import (
	"context"
	"sync"
)

type Pool struct {
	Workers int
	Jobs    chan Job
	Results chan Result
}

func NewPool(workers, queueSize int) *Pool {
	return &Pool{
		Workers: workers,
		Jobs:    make(chan Job, queueSize),
		Results: make(chan Result, queueSize),
	}
}

// Run starts workers and waits for them to finish.
func (p *Pool) Run(ctx context.Context, totalJobs int, sink func(Result)) {

	// 1. Start Workers
	var wg sync.WaitGroup
	for i := 0; i < p.Workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, p.Jobs, p.Results)
		}(i)
	}

	// 2. Produce jobs
	go func() {
		for i := 0; i < totalJobs; i++ {
			select {
			case <-ctx.Done():
				return
			case p.Jobs <- Job{ID: i, Payload: "task"}:
			}
		}
		close(p.Jobs)
	}()

	// 3. Collect results in a separate gorountine
	done := make(chan struct{})
	go func() {
		for r := range p.Results {
			sink(r)
			if r.JobID == totalJobs-1 {
				// last job processed: user may stop early
			}
		}
		close(done)
	}()

	// 4. Wait for Workers to exit
	wg.Wait()
	close(p.Results)
	<-done
}
