# Code facts — redis-performance/redis-benchmark-connect

There is no reviewer-comment precedent in this repo (see `review-history.md`), so
this file substitutes the only other honest source of grounding: what the actual
codebase looks like today, read directly rather than inferred. Use these as things
worth checking on a real diff — not as "maintainers have said this matters," since
no maintainer here has ever said anything in review. Re-verify against the current
diff/files; this snapshot is from 2026-08-26/27 and the repo is small enough to
re-read in full if something looks off.

## What this project is

A single-purpose Go CLI (`main.go`, ~267 lines, package `main`, one file, no
internal packages) that opens N TCP or TLS connections to a Redis server
(sequentially or with one goroutine per connection) and reports total elapsed time
and average per-connection latency, optionally appending per-connection timings to
`output.txt`. It depends on `github.com/gomodule/redigo/redis` for the actual Redis
protocol dial. There is no server-side code, no persistence layer, no HTTP surface —
review scope is almost always going to be `main.go`, `go.mod`/`go.sum`,
`.github/workflows/ci.yml`, or `tests/*.sh`.

## Testing reality — read CONTRIBUTING.md's rule against what actually exists

`CONTRIBUTING.md` states: "All new behaviour must be covered by tests... Coverage
should not decrease." As of this mining, **there are zero Go unit tests in this
repo** — no `*_test.go` file exists anywhere. The only tests are two shell scripts
(`tests/run_tests.sh`, `tests/gen-test-certs.sh`) that build the binary, start a real
(TLS-enabled) Redis, run the compiled tool against it, and check the exit code —
end-to-end smoke tests, not unit tests, and there is no code-coverage tool wired
into `.github/workflows/ci.yml` at all. So: the written rule exists, but there is
currently no mechanism, human or automated, that would actually catch "new behavior
shipped with zero test coverage" the way Codecov does on a repo like redisbench-admin.
For a new CLI flag or new behavior, checking whether `tests/run_tests.sh` (or a new
shell test) actually exercises it is a real, meaningful thing to look for — but be
accurate that failing to do so isn't currently enforced by any tool, only by
CONTRIBUTING.md's text.

## Concrete patterns already present in `main.go` worth knowing before reviewing a diff near them

These are observations about the current code, not bugs necessarily worth
re-litigating on every PR — they exist so that if a PR touches these exact spots,
you're reviewing with real context instead of a blank slate:

- **The parallel path accumulates into a shared, unsynchronized variable.**
  `testAndMeasureConnectionsParallel` launches one goroutine per connection, and each
  goroutine does `totalConnectionTime += float32(connElapsedTime.Milliseconds())` on
  a `float32` declared once outside the loop, with no mutex/atomic and no
  per-goroutine accumulation-then-sum. This is a real, currently-shipping data race
  under Go's race detector (not something CI currently runs `-race` to catch — the
  CI workflow builds and runs the binary directly, not `go test -race`). If a PR
  touches this function, this existing race is fair, concrete, on-point context —
  and a new PR that adds more shared mutable state to the parallel path without
  synchronization would be repeating a pattern already latent here, worth naming
  specifically rather than a generic "watch out for concurrency" comment.
- **`output.txt` is opened and `defer`-closed inside the per-connection loop/goroutine
  body**, in both the sequential and parallel functions — i.e. the file handle from
  iteration 1 isn't closed until the *enclosing function* returns (sequential case)
  or until that goroutine returns (parallel case, so N file handles briefly overlap
  under `-parallel`), not closed after each write. It works today because
  `numConnections` is bounded and short-lived, but a PR that raises the practical
  scale of `-numConnections` materially, or restructures this loop, is worth checking
  for fd exhaustion.
- **TLS auto-enable via `flag.Lookup`.** After `flag.BoolVar(&useTLS, "tls", ...)`
  already binds `useTLS` directly, `main()` separately does
  `if flag.Lookup("tls") != nil && flag.Lookup("tls").DefValue !=
  flag.Lookup("tls").Value.String() { useTLS = true }` — functionally redundant with
  the direct binding (both approaches derive `useTLS` from the same flag), just via a
  second, more roundabout mechanism. Not a bug, but if a PR touches this area, it's
  worth asking why the second path exists rather than assuming it's load-bearing.
- **The `hello` path sleeps 10 full seconds per connection**
  (`if hello { conn.Do("HELLO"); time.Sleep(10 * time.Second); ... }`) in both the
  sequential and parallel functions, executed *per connection* — with `-parallel` off
  and `-numConnections` above a handful, `-hello` makes a run take
  `10s * numConnections` well before the connection overhead this tool exists to
  measure is even relevant. If a PR touches the `-hello` path, worth asking whether
  that sleep is intentional (e.g. waiting out something server-side) and, if so,
  whether it should be a flag rather than a hardcoded constant — there's no commit
  message or comment in the current code explaining why 10s specifically.
- **The version string is a hardcoded Go var**: `var version = "1.0.3"` at the top of
  `main.go`, printed by `-version`. A version-bump PR needs to touch this literal;
  there's no build-time ldflags injection or git-tag-derived version here.

## Written rules actually in this repo (apply these; they are real, not mined-from-comments)

From `AGENTS.md` (added by PR#3, still current):
- Branch naming `<type>/<short-description>`.
- "Do not add comments that describe *what* the code does — only add comments when
  the *why* is non-obvious" — several existing comments in `createTLSConfig`
  (the `InsecureSkipVerify` block) already follow this; hold new code to the same bar.
- "Do not introduce new dependencies without checking with the maintainer" — this
  project currently has exactly one non-stdlib dependency (`redigo`); a new import
  of anything beyond the Go standard library is worth flagging explicitly for that
  reason alone.
- Do not reformat unrelated files; do not remove error handling or tests; do not
  commit secrets/credentials/large binaries; do not amend published commits.

From `CONTRIBUTING.md`:
- One logical change per PR.
- Follow existing conventions (formatting, naming, error handling) — this codebase's
  existing error handling is uniformly "print with `fmt.Println`/`fmt.Printf` and
  `return`," never `log.Fatal`, panics, or wrapped/sentinel errors; match that unless
  there's a stated reason not to.
- No dead code, no commented-out blocks.
- Coverage should not decrease (see the testing-reality note above for how literally
  this is currently enforceable).

## Scope note

If a PR touches something entirely outside this surface (e.g. an unrelated vendored
asset, a totally different subsystem this repo doesn't have), say so in one sentence
and treat it as out of scope rather than force-fitting this checklist.
