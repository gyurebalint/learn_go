# High-Performance Log Ingester (CLI)

A concurrent, fault-tolerant ETL pipeline written in Go. This tool streams large JSON datasets, processes them in parallel using a Worker Pool pattern, and safely ingests them into a SQLite database using Write-Ahead Logging (WAL).

## 🚀 Key Features

* **Streaming Architecture:** Uses `io.Reader` and `json.Decoder` to process datasets of arbitrary size (GBs/TBs) with constant O(1) memory usage.
* **Concurrency:** Implements a Fan-Out/Fan-In pattern with a configurable Worker Pool.
* **Resilience:** Full `context.Context` propagation for timeout management and cancellation.
* **Graceful Shutdown:** Handles OS signals (`SIGINT`, `SIGTERM`) to ensure database transactions complete before exiting.
* **Database:** High-throughput SQLite implementation using WAL (Write-Ahead Logging) mode and exponential backoff strategies (handled via driver config).

## 🛠 Project Structure

Adheres to the Standard Go Project Layout:

```text
log-ingester/
├── cmd/
│   └── ingester/    # Application entry point (wiring)
├── internal/
│   ├── models/      # Domain entities (Person)
│   └── storage/     # Database implementations (Repository pattern)
├── users.json       # Sample dataset
├── go.mod           # Dependency definitions
└── Makefile         # Build scripts