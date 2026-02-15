# Order Processor

A concurrent order processing system in Go that maintains **user-specific queues** to ensure sequential order processing per user while allowing parallel processing across different users.

This project demonstrates:

- Thread-safe storage for user balances.
- Worker pool for concurrent task execution.
- Safe shutdown and order queue management.
- End-to-end testing and unit testing with mocks.

---

## Features

- **OrderProcessor API**: Submit orders, query user balances, gracefully shutdown.
- **Storage**: In-memory thread-safe key-value storage mapping user IDs to balances.
- **WorkerPool**: Concurrent processing of tasks with configurable worker count and queue buffer.
- **User-specific queues**: Ensures orders from the same user are processed sequentially.
- **Graceful shutdown**: Waits for all tasks to complete before closing workers.

---

## Getting Started

### Requirements

- Go 1.25+
- [`just`](https://github.com/casey/just) – Command runner for running tasks
- [`Ginkgo v2`](https://github.com/onsi/ginkgo) – BDD-style testing framework
- [`Gomega`](https://github.com/onsi/gomega) – Matcher library for tests
- [`GoMock`](https://github.com/uber-go/mock) – For generating mocks

### Installation

Clone the repository:

```bash
git clone https://github.com/antoniuk-oleksandr/order_processor.git
cd order_processor
```

Install dependencies:

```bash
go mod tidy
```

## Usage

### Running the Example

```bash
just run
# or
go run ./cmd/order_processor
```

This runs the main program that submits sample orders and prints user balances.

## Testing

Run all tests:

```bash
just test
# or
ginkgo -r -race -v
```

Run end-to-end tests only:

```bash
just e2e
# or
ginkgo -r -v -race --label-filter="e2e"
```

Run unit tests only:

```bash
just unit
# or
ginkgo -r -v --label-filter="unit"
```

Tests use:

- Ginkgo/Gomega for expressive BDD-style testing.
- GoMock for mocking storage and worker pool interfaces.

## Documentation

API documentation is available via GoDoc:

```bash
just doc
# or
pkgsite -open .
```


## Justfile Commands

- `just run` – Run the main program.
- `just doc` – Generate and open API documentation.
- `just cov` – Generate test coverage report.
- `just test` – Run all tests.
- `just e2e` – Run end-to-end tests.
- `just unit` – Run unit tests only.