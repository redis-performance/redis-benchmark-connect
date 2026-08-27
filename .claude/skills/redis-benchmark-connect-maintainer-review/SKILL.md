---
name: redis-benchmark-connect-maintainer-review
description: Review a redis-performance/redis-benchmark-connect pull request, branch, or diff, grounded in this repo's actual (very thin) GitHub history and its real AGENTS.md/CONTRIBUTING.md rules rather than generic Go code-review advice or a borrowed "maintainer personality." Use this whenever the user asks to review a redis-benchmark-connect PR "like a maintainer would," asks whether a redis-benchmark-connect PR would pass real review, wants a repo-specific pre-merge check, or is deciding accept/reject on a redis-performance/redis-benchmark-connect PR. Prefer this over a generic code-review skill for this repo — a generic skill won't know that this repo has zero written review comments in its entire history, or the specific concurrency/resource-handling patterns already present in its one source file.
---

# redis-benchmark-connect maintainer-style review

**Read `references/review-history.md` and `references/code-facts.md` before writing
anything.** This skill exists to keep a review grounded in what is actually true of
this repo, and the single most important thing that's true of it is this:

## The central fact: there is no review-comment history to imitate

`redis-benchmark-connect` has had exactly five pull requests, ever, all merged, and
**not one of them carries a single written review comment** — no top-level review
body, no inline comment, no issue comment. Two PRs have a bare `APPROVED` from
`paulorsousa` with an empty body; the rest have no recorded review at all. Zero
issues have ever been filed against this repo. Full detail, including every PR and
who touched it, is in `review-history.md`.

This is meaningfully thinner than even a small repo's history usually is — do not
treat this skill as a smaller version of a skill for a repo with real (if sparse)
review quotes to draw on. There is no "voice" to reconstruct here, because no one has
ever written one down. **Say this plainly if it's relevant, and do not invent a
maintainer personality, a recurring nitpick pattern, or a quoted precedent that
doesn't exist.** If you catch yourself writing something that sounds like "this
project's maintainers usually flag X" or "as seen in past reviews," stop — check
`review-history.md`, and if there's no citation for it there, don't say it.

## What to ground a review in instead

With no comment history to mine, ground the review in the two things that actually
are real:

1. **The written rules in `AGENTS.md` and `CONTRIBUTING.md`** (both already in the
   repo, added in PR#3) — branch naming, no unrelated reformatting, no new
   dependencies without asking, "why not what" comments only, one logical change per
   PR, and the stated (if not currently tool-enforced) testing/coverage expectation.
2. **The actual code**, read fresh for the diff at hand, using `code-facts.md`'s
   documented properties of the current `main.go`/tests/CI as context — e.g. if a PR
   touches the parallel connection path, you now know there's already an
   unsynchronized shared accumulator there worth being precise about; if it touches
   `-hello`, you know about the existing 10-second per-connection sleep. Don't
   silently assume these facts are still true forever — this is a small, actively
   read file; re-check `main.go` if something in `code-facts.md` looks stale against
   the current diff.

## Process

1. **Get the material.** `gh pr view <n> --repo redis-performance/redis-benchmark-connect --json body,commits,files,author` and `gh pr diff <n> --repo redis-performance/redis-benchmark-connect`. Read the full PR description — this project's real author (`fcostaoliveira`/`filipecosta90`) has written real "Root cause" / "Changes" sections on recent PRs (#4, #5); if the description already explains something, acknowledge that instead of re-deriving it.

2. **Scope check.** This is a single-file Go CLI plus two shell test scripts and one CI workflow — see `code-facts.md` for the full inventory. If the diff falls entirely outside that surface, say so in one sentence and treat it as out of scope rather than force-fitting a checklist built for a Go connection-benchmarking tool.

3. **Read the diff against `code-facts.md` and the written rules**, not against a memory of imagined past reviews:
   - Does new behavior have any test coverage at all, given there is currently no unit-test infrastructure and only shell-based end-to-end smoke tests (`tests/run_tests.sh`)? Naming this is legitimate per `CONTRIBUTING.md`'s explicit rule, but don't claim it will mechanically block merge — nothing in this repo's CI currently enforces it.
   - Does the diff add a new dependency? `AGENTS.md` explicitly asks that this be checked with the maintainer first, and this project currently has exactly one (`redigo`) — worth naming explicitly, not glossing over.
   - Does the diff touch the parallel-connections path, the per-connection file-output path, or the `-hello` path? If so, read `code-facts.md`'s notes on the existing unsynchronized accumulator, the loop-scoped `defer file.Close()`, and the hardcoded 10s sleep, and check whether the new code repeats, worsens, or (ideally) fixes the same pattern — cite the actual line/mechanism, not a generic "watch for race conditions."
   - Does new code match this codebase's actual existing conventions — error handling via `fmt.Println`/`return`, not panics/`log.Fatal`; comments that explain *why*, not *what* (per `AGENTS.md`, verbatim)?
   - Is a hardcoded literal like `version = "1.0.3"` bumped where relevant?

4. **Calibrate how much to say by what's actually there, biased toward this repo's own real bias toward silence.** A PR that is docs-only, a dependency/version bump, or a small change with no correctness/resource-handling implications is exactly the kind of PR this repo's real history shows getting merged with zero comment (three of five real PRs here got none at all, including a 236-line feature PR) — `skip_comment=true` is the *faithful* choice for those, not a cop-out. Reserve an actual comment for something concrete: a real gap against a written rule, or a real, specific mechanism in the diff worth naming (per step 3), not commentary manufactured to look thorough.

5. **Write briefly, in plain prose, if you write anything at all.** A few sentences, specific and concrete (name the function, the flag, the exact behavior), never a formal "Correctness / Security / Performance" essay — nothing in this project's real history looks like that, and a wall of generic Go-idiom advice would be less faithful to this repo than saying nothing. Never literally `@`-mention a GitHub username, even to suggest a second opinion — say so in prose instead ("worth a second look from whoever's touched the TLS path before").

6. **Land on a verdict that matches what actually happens here**: silence (`skip_comment=true`) is the modal real outcome; when you do comment, a short, concrete note is enough — this repo has no real precedent for `APPROVED`-with-substantive-comment or `COMMENTED`-with-a-blocking-ask, so don't manufacture false confidence about which GitHub review state a real maintainer would pick. Never write the word "Verdict" or a bolded/labeled summary line; end in plain prose.

## What NOT to do

- Don't borrow the memtier_benchmark or redisbench-admin skills' reviewer quotes,
  voice profiles, or nitpick taxonomies — those are real precedent for *different*
  repos with real review history; none of it is evidenced here.
- Don't invent a "maintainer usually says X" claim not traceable to
  `review-history.md` or a written rule in `AGENTS.md`/`CONTRIBUTING.md`.
- Don't claim the testing/coverage rule in `CONTRIBUTING.md` is mechanically
  enforced — it currently isn't (see `code-facts.md`).
- Don't manufacture a lengthy review to seem thorough on a routine PR — real history
  here overwhelmingly favors silence; match it.
- Don't literally `@`-mention any GitHub username, ever.
- Don't close with a labeled, bolded verdict block — end in plain prose.
