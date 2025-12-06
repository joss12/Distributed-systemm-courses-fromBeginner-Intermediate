package pinger

import (
	"context"
	"net/http"
	"time"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// PingOne issues a GET with a pre-request timeout and measure latency.
func PingOne(ctx context.Context, client Doer, url string, timeout time.Duration) Result {
	start := time.Now()

	//Bind timeout to the request's context; cancellation propagates to transport.
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Timestamp: time.Now(), Error: err.Error()}
	}

	//You could costumize http.client (timeouts, transports). We accept it as Dor.
	resp, err := client.Do(req)
	lat := time.Since(start)

	if err != nil {
		return Result{
			URL:       url,
			Latency:   lat,
			Timestamp: time.Now(),
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()

	return Result{
		URL:       url,
		Status:    resp.StatusCode,
		Latency:   lat,
		Timestamp: time.Now(),
	}

}

func PingWithRetry(ctx context.Context, client Doer, url string, timeout time.Duration) Result {
	attempts := 3
	backoff := 100 * time.Microsecond

	for i := 1; i <= attempts; i++ {
		r := PingOne(ctx, client, url, timeout)
		if r.Error == "" {
			return r
		}

		//if this was the last attempts, return it
		if i == attempts {
			return r
		}

		//respect context cancellation
		select {
		case <-ctx.Done():
			return r
		case <-time.After(backoff):
			backoff *= 2 //exponential backoff
		}
	}
	return Result{} //unreachable
}

func NewDefaultclient() *http.Client {
	return &http.Client{
		Timeout: 0, //we prefer pre-request Context timeout, not client-wide
		// Transport: &http.Transport{...} //(left default for simplicity)
	}
}
