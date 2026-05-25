# Contributing

We treat this repo as "Open Source" within Redis: anyone who clears the bar below is welcome to contribute.

## Local setup

```bash
git clone git@github.com:redis-performance/redis-benchmark-connect.git
cd redis-benchmark-connect
go mod download
go build -o redis_benchmark_connect .
```

Requires **Go 1.20 or later**.

For TLS-related testing you also need `openssl` available on your PATH.

## Branch naming

```
<type>/<short-description>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

Example: `feat/add-pipeline-mode`

## Coding standards

- Keep changes focused; one logical change per PR.
- Follow the conventions already present in the codebase (formatting, naming, error handling).
- No dead code, no commented-out blocks.

## Submitting changes

1. Fork or create a branch from `main`.
2. Make your changes with clear, atomic commits.
3. Open a pull request against `main` with a descriptive title and summary.
4. Address review comments promptly; force-push to the same branch to update.

## Testing

- All new behaviour must be covered by tests.
- Existing tests must pass: run the test suite locally before opening a PR.
- Coverage should not decrease.

Run the full test suite (requires a local Redis instance on port 6379):

```bash
# Build the binary first
go build -o redis_benchmark_connect .

# Run plain TCP tests
./tests/run_tests.sh

# Run TLS tests (generates self-signed certs via openssl, starts Redis on port 6380)
./tests/gen-test-certs.sh
TLS=1 ./tests/run_tests.sh
```

Individual TLS protocol variants:

```bash
TLS_PROTOCOLS="tlsv1.2" TLS=1 ./tests/run_tests.sh
TLS_PROTOCOLS="tlsv1.3" TLS=1 ./tests/run_tests.sh
```

CI runs the same steps across Ubuntu on Go 1.20.x and 1.21.x (see `.github/workflows/ci.yml`).

## Review process

- At least one maintainer approval is required before merge.
- CI must be green.
- Maintainers may request changes or close PRs that do not meet the bar — this is normal and not personal.