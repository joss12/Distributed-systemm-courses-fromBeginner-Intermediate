package pinger

import (
	"context"
	"fmt"
)

func RunFanOutFanIn(ctx context.Context, urls []string, cfg Config, fn func(Result)) {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 0
	}

	sem := make(chan struct{}, cfg.Concurrency) //semaphore
	results := make(chan Result)

	//producer: for each URL, spawn a worker respecting Concurrency
	go func() {
		for _, u := range urls {
			select {
			case <-ctx.Done():
				break
			case sem <- struct{}{}: //acquire "slot"
				// spawn worker
				go func(url string) {
					fmt.Printf("[worker-start]%s\n", url)

					defer func() {
						fmt.Printf("[worker-end]%s\n", url)
						<-sem
					}() //release slot
					//r := PigOne(ctx, NewDefaultClient(), url, cfg.Timeout)
					r := PingWithRetry(ctx, NewDefaultClient(), url, cfg.Timeout)
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
		fn(r)
	}
}
