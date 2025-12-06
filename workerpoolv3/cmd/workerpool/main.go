package main

import (
	"context"
	"flag"
	"time"

	"github.com/workerpoolv3/internal/pool"
)

func main() {
	var (
		workers   int
		totalJobs int
		queueSize int
	)

	flag.IntVar(&workers, "workers", 4, "number of workers")
	flag.IntVar(&totalJobs, "jobs", 20, "total jobs")
	flag.IntVar(&queueSize, "queue", 20, "queue buffer size")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p := pool.NewPool(workers, queueSize)

	p.Run(ctx, totalJobs, pool.PrintSink)
}
