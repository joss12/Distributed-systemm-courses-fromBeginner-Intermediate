
# Chapter 1 — Distributed Pinger (Fan-Out / Fan-In)

A high-performance HTTP pinger written in **Go**, demonstrating the core
distributed-systems pattern **fan-out / fan-in**:

- Fan-out: spawn bounded concurrent workers  
- Fan-in: aggregate results safely  
- Per-request timeouts using `context`  
- Concurrency limits using a semaphore  
- Optional JSON & CSV export  
- Optional URL input file  
- Retry with exponential backoff  

This chapter forms the foundation for all later distributed systems work:
replication, consensus, scheduling, heartbeats, leader election, etc.

---

# 🚀 Features

- 🔥 **Bounded concurrency** (avoid overwhelming remote hosts)
- ⏱ **Per-request timeout** (context cancellation)
- 🔁 **Retry with exponential backoff**
- 📊 **Real-time logging**
- 📁 **Load URLs from file**
- 📦 **CSV & JSON export**
- 🧱 **Clean cmd/internal architecture**

---

# 🧠 Architecture

```
                       +---------------------+
URLs  ───────────▶     |  Runner (fan-out)   |
                       |  - concurrency sem  |
                       |  - spawn workers    |
                       +----------┬----------+
                                  │
                                  ▼
                        +------------------+
                        |   Ping Worker    |
                        |  - HTTP GET      |
                        |  - Timeout       |
                        |  - Retries       |
                        +--------┬---------+
                                 │
                                 ▼
                       +----------------------+
                       |  Fan-In Aggregation  |
                       |   Callback (sink)    |
                       +----------------------+
```

---

# 📁 File Tree

```
pinger/
├─ cmd/
│  └─ pinger/
│     └─ main.go
├─ internal/
│  └─ pinger/
│     ├─ model.go
│     ├─ httpclient.go
│     ├─ runner.go
│     └─ csv.go
└─ go.mod
```

---

# 📦 Installation

```bash
git clone <your-repo>
cd pinger
go mod tidy
```

---

# 🏃 Usage

## Basic usage

```bash
go run ./cmd/pinger --concurrency 8 --timeout 1200ms https://example.com https://httpbin.org/delay/2
```

## Load URLs from a file

```
urls.txt:
https://example.com
https://google.com
https://httpbin.org/status/200
https://httpbin.org/delay/1
```

```bash
go run ./cmd/pinger --urls-file urls.txt --concurrency 8 --timeout 1200ms
```

## Export results to CSV

```bash
go run ./cmd/pinger --csv results.csv https://example.com https://httpbin.org/delay/1
```

## Export results to JSON

```bash
go run ./cmd/pinger --json results.json https://example.com https://httpbin.org/delay/1
```

## CSV + JSON + URL file

```bash
go run ./cmd/pinger --urls-file urls.txt --csv out.csv --json out.json --timeout 1s --concurrency 4
```

---

# 🧪 Implemented Exercises

1. **Concurrency visualization** → `internal/pinger/runner.go`  
2. **Load URLs from file (`--urls-file`)** → `cmd/pinger/main.go`  
3. **Retry with exponential backoff** → `internal/pinger/httpclient.go`  
4. **JSON export (`--json`)** → `cmd/pinger/main.go`  

---

# 📄 Sample CSV Output

```
timestamp,url,status,latency_ms,error
2025-01-10T14:12:03.001Z,https://example.com,200,87.231,
2025-01-10T14:12:03.100Z,https://httpbin.org/delay/2,0,1200.562,context deadline exceeded
```

---

# 🧩 Next Steps

- Chapter 2 — Single Node Key-Value Store  
- Chapter 3 — Worker Pool  
- Chapter 4 — Distributed KV Store  
- Chapter 5 — Leader Election (Raft-lite)