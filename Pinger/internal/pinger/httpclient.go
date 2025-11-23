package pinger

import (
	"context"
	"net/http"
	"time"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// PigOne issues a GET with a per-request timeout and measures latency.
func PingOne(ctx context.Context, client Doer, url string, timeout time.Duration) Result {
	start := time.Now()

	//Bind timeout to the request's context; cancellation progagates top transport.
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Timestamp: time.Now(), Error: err.Error()}
	}

	//You could costomize http.Client (timeouts, transports). We accept it as Doer.
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

// PingWithRetry retries PingOne using exponential backoff.
func PingWithRetry(ctx context.Context, client Doer, url string, timeout time.Duration) Result {
	attempts := 3
	backoff := 100 * time.Millisecond

	for i := 1; i <= attempts; i++ {
		r := PingOne(ctx, client, url, timeout)
		if r.Error == "" {
			return r
		}

		//if this was the last attempt, return it
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
	return Result{} // unreachable
}

// NewDefaultClient returns a sane http.Client for CLI tools.
func NewDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 0, //we prefer per-request Context timeouts, not client-wide
		// Transport: &http.Transport{ ... } // (left default for simplicity)
	}
}
