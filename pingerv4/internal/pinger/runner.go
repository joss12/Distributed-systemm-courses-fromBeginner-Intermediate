package pinger

import (
	"context"
	"fmt"
)

// RunFanOutFanIn concurrent pings with a concurrent limit.
// - It never spawn more than cfg.Concurrent goroutines at once.
// - It returns only after all URLs are processed.
func RunFanOutFanIn(ctx context.Context, urls []string, cfg Config, fn func(Result)) {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 0
	}

	sem := make(chan struct{}, cfg.Concurrency) //semaphore to bound Concurrency
	results := make(chan Result)

	//Producer: for each URL, spawn a worker respecting concurrency
	go func() {
		for _, u := range urls {
			select {
			case <-ctx.Done():
				break
			case sem <- struct{}{}: //acquire "slot"
				//spawn worker
				go func(url string) {
					fmt.Printf("[Worker-start]%s\n", url)

					defer func() {
						fmt.Printf("[Worker-end]%s\n", url)
						<-sem
					}() //release slot
					//r := PingOne(ctx, NewDefaultclient(), url, cfg.Timeout)
					r := PingWithRetry(ctx, NewDefaultclient(), url, cfg.Timeout)
					results <- r
				}(u)
			}
		}
		//Wait for all workers to release slots: drain the semaphore fully
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
		close(results)
	}()

	//consumer: fan-in
	for r := range results {
		fn(r) //callemr-defined sink(print, collect, export)
	}
}
