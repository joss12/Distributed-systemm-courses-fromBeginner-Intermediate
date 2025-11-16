package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/concurrent-downloader/internal/pinger"
)

func main() {
	var (
		timeoutStr   string
		concurrency  int
		writeCSVPath string
		jsonPath     string
		urlsFile     string
	)

	flag.StringVar(&timeoutStr, "timeout", "1200ms", "per-request timeout (e.g. 8000ms, 2s)")
	flag.IntVar(&concurrency, "concurrency", 8, "max in-flight requests")
	flag.StringVar(&writeCSVPath, "cvs", "", "optional path to write results csv")
	flag.StringVar(&jsonPath, "json", "", "write results to JSON file")
	flag.Parse()

	//if flag.NArg() == 0 {
	//	fmt.Fprintf(os.Stderr, "usage: %s [flags] <url1> <url2> ...\n", os.Args[0])
	//	flag.PrintDefaults()
	//	os.Exit(2)
	//}

	var urls []string
	if urlsFile != "" {
		data, err := os.ReadFile(urlsFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read urls-file:", err)
			os.Exit(2)
		}
		urls = strings.Split(strings.TrimSpace(string(data)), "\n")
	} else {
		if flag.NArg() == 0 {
			fmt.Fprintf(os.Stderr, "usage: %s [flags] <url1> <url2>...\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(2)
		}
		urls = flag.Args()
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid -timeout:", err)
		os.Exit(2)
	}

	cfg := pinger.Config{
		Timeout:     timeout,
		Concurrency: concurrency,
	}
	ctx := context.Background()

	//Collect results (and also print them line-by-line).
	var collected []pinger.Result
	sink := func(r pinger.Result) {
		status := r.Status
		if r.Error != "" {
			fmt.Printf("%-40s error=%-24s latency=%v\n", r.URL, truncate(r.Error, 24), r.Latency)
		} else {
			fmt.Printf("%-40s status=%-3d        latency=%v\n", r.URL, status, r.Latency)
		}
		collected = append(collected, r)
	}

	//pinger.RunFanOutFanIn(ctx, flag.Args(), cfg, sink)
	pinger.RunFanOutFanIn(ctx, urls, cfg, sink)

	if writeCSVPath != "" {
		f, err := os.Create(writeCSVPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "csv create:", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pinger.WriteCSV(f, collected); err != nil {
			fmt.Fprintln(os.Stderr, "cvs write", err)
			os.Exit(1)
		}
		fmt.Println("cvs written to", writeCSVPath)
	}

	if jsonPath != "" {
		f, err := os.Create(jsonPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "json create:", err)
			os.Exit(1)
		}
		defer f.Close()

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(collected); err != nil {
			fmt.Fprintln(os.Stderr, "json write:", err)
			os.Exit(1)
		}
		fmt.Println("json written to", jsonPath)
	}

}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	//rough cut that tries to preserve the right side detail (like "i/o timeout")
	if n <= 3 {
		return s[:n]
	}
	left := n - 3
	return s[:left]
}
