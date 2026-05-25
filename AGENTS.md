# Agent guidelines

Instructions for AI coding agents (Claude Code, Copilot, Cursor, etc.) working in this repo.

## Project overview

`redis-benchmark-connect` is a Go CLI tool that benchmarks the time it takes to establish connections to a Redis server. It opens a configurable number of TCP or TLS connections (sequentially or in parallel), reports total elapsed time and average connection latency per connection, and optionally writes per-connection timings to a file. It is useful for measuring Redis connection overhead under different network and TLS configurations.

## Local setup

```bash
git clone git@github.com:redis-performance/redis-benchmark-connect.git
cd redis-benchmark-connect
go mod download
go build -o redis_benchmark_connect .
```

Requires **Go 1.20 or later** and `openssl` on the PATH for TLS certificate generation.

## Branch naming

Same as human contributors: `<type>/<short-description>` (e.g. `fix/off-by-one-in-pipeline`).

## Coding standards

- Match the style already in the file you are editing.
- Prefer clear, minimal changes over large refactors unless explicitly asked.
- Do not add comments that describe *what* the code does — only add comments when the *why* is non-obvious.
- Do not introduce new dependencies without checking with the maintainer.

## Running tests

Build the binary, then run the shell-based integration tests:

```bash
# Build
go build -o redis_benchmark_connect .

# Plain TCP (requires Redis on localhost:6379)
./tests/run_tests.sh

# TLS (generates certs, requires Redis supporting TLS on port 6380)
./tests/gen-test-certs.sh
TLS=1 ./tests/run_tests.sh

# Specific TLS protocol versions
TLS_PROTOCOLS="tlsv1.2" TLS=1 ./tests/run_tests.sh
TLS_PROTOCOLS="tlsv1.3" TLS=1 ./tests/run_tests.sh
```

Always run tests before declaring a task complete.

## How to submit changes

1. Create a branch: `git checkout -b <type>/<description>`.
2. Commit with a clear message focused on *why*, not *what*.
3. Open a pull request against `main`.
4. Do **not** push directly to `main`.

## What to avoid

- Do not reformat files unrelated to your change.
- Do not remove error handling or tests.
- Do not commit secrets, credentials, or large binary files.
- Do not amend published commits.