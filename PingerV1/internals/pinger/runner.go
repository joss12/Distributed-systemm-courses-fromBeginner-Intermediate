package pinger

import (
	"context"
	"fmt"
)

func RunFanOutFanIn(ctx context.Context, urls []string, cfg Config, fn func(Result)) {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 0
	}

	sem := make(chan struct{}, cfg.Concurrency) //semaphore to bound Concurrency
	results := make(chan Result)

	//producer: for ear URL, spawn a worker respecting concurrency,
	go func() {
		for _, u := range urls {
			select {
			case <-ctx.Done():
				break
			case sem <- struct{}{}:
				//spawn worker
				go func(url string) {
					fmt.Printf("[work-start]%s\n", url)

					defer func() {
						fmt.Printf("[worker-end] %s\n", url)
						<-sem
					}() //replace slot
					//r := PingOne(ctx, NewDefaultClient(), url, cfg.Timeout)
					r := PingWithRetry(ctx, NewDefaultClient(), url, cfg.Timeout)
					results <- r
				}(u)
			}
		}
		//Wait for all worjkers to release slots: drain the semaphore fully
		for i := 0; i < cap(sem); i++ {
			sem <- struct{}{}
		}
		close(results)
	}()

	//consumer:fan-in
	for r := range results {
		fn(r)
	}
}
