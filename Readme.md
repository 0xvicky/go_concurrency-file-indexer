# Concurrent File Indexer (Go)

A high-performance concurrent file indexing tool written in Go.
The project explores core systems-programming ideas such as **goroutines, worker pools, channels, synchronization, and filesystem I/O performance**.

The program recursively scans a directory, distributes file paths through a channel-based task queue, and processes them using a pool of workers that compute **SHA256 hashes** of file contents. Results are stored in a thread-safe shared map.

The goal of the project is educational: understanding how Go enables efficient **concurrent pipelines for I/O-heavy workloads**.

---

# Architecture

The system follows a **producer–consumer pipeline model**.

Filesystem Scanner (Producer) → Channel (Task Queue) → Worker Pool → Hash + Storage

1. **Scanner**

   * Recursively walks the filesystem
   * Sends discovered file paths into a channel

2. **Channel (Queue)**

   * Acts as a thread-safe queue
   * Decouples file discovery from file processing

3. **Worker Pool**

   * Multiple goroutines consume tasks from the queue
   * Each worker reads the file, computes SHA256, and stores the result

4. **Storage**

   * A mutex-protected map stores `filePath → hash`

5. **Lifecycle Management**

   * `WaitGroup` ensures workers exit cleanly
   * Channel is closed once scanning completes

---

# Key Concepts Demonstrated

### Goroutines

Lightweight threads used to process files concurrently.

### Channels

Used as a task queue between producer and consumers.

### Worker Pool Pattern

Controls concurrency and prevents uncontrolled goroutine creation.

### Synchronization

A mutex protects shared storage from race conditions.

### Pipeline Architecture

Work is divided into stages: **scan → queue → process → store**.

---

# Performance Observation

Benchmarking revealed an important systems insight.

First run:
~21 seconds

Subsequent runs:
~140 ms

The reason is **filesystem page caching**.

During the first execution the operating system must read files from disk. Afterward, the OS caches those files in memory. Subsequent runs read from RAM instead of disk, dramatically reducing runtime.

This demonstrates the major performance gap between:

Disk I/O → milliseconds
Memory access → nanoseconds

---

# Project Structure

```
concurrent-indexer
│
├── cmd/
│   └── indexer/
│       └── main.go
│
├── internal/
│   ├── scanner/      # Filesystem traversal
│   ├── worker/       # Worker pool logic
│   ├── hashing/      # SHA256 hashing
│   ├── storage/      # Thread-safe result storage
│   └── config/       # Shared constants and configuration
│
├── go.mod
└── README.md
```

### cmd/

Contains application entrypoints (binaries).

### internal/

Contains the core system components.

Each package handles a single responsibility, following a modular and production-style Go layout.

---

# Running the Program

Clone the repository:

```
git clone <repo-url>
cd concurrent-indexer
```

Run the indexer:

```
go run ./cmd/indexer
```

Or build a binary:

```
go build ./cmd/indexer
./indexer
```

---

# Lessons Learned

1. **Concurrency hides I/O latency** but does not eliminate it.
2. **Printing inside tight loops severely degrades performance** due to stdout I/O.
3. **Operating system caching dramatically affects benchmarks.**
4. Proper **package boundaries improve maintainability**.
5. Worker pools are a fundamental pattern in systems engineering.

---

# Future Improvements

Possible extensions to evolve the project into a more advanced system:

• Replace in-memory storage with **SQLite or BoltDB**
• Implement **content-addressable storage** (deduplicated files)
• Add **progress metrics and observability**
• Build a **distributed indexing system** across multiple nodes
• Implement **parallel directory scanning**

---

# Technologies Used

Go
Goroutines
Channels
SHA256 hashing
Filesystem APIs

---

# Purpose

This project is part of a systems-programming learning journey focused on building high-performance backend and infrastructure tools using Go.

It demonstrates how simple primitives—channels, goroutines, and synchronization—can form scalable concurrent systems.
