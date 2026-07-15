# Demiurge — Agent Guide

## What this project is

A macOS CLI tool (`demi`) for developers who maintain many projects simultaneously. It provides terminal-native project awareness — not a project manager, no tasks or boards. Local-only, no accounts, no cloud. See `README.md` for full context.

## Before starting any work

Run this to get current state:

```sh
gh issue list --state open
gh api repos/{owner}/{repo}/milestones
gh pr list
```

**Every piece of work must be linked to at least one open issue.** If there is no issue that covers what you are about to do, stop and ask the user whether to create one before proceeding. Do not open issues yourself without asking first.

## Issue workflow

- All work (features, fixes, chores that change behaviour) links to an issue
- Reference issues in commit messages: `feat: add git status view, closes #7`
- Multiple related commits can reference the same issue without closing it: `ref #7`
- Issues close automatically when the commit lands on `main` if `closes #N` appears in it

## Sub-issues

Use sub-issues to break down anything non-trivial. A parent issue describes the feature or goal at a high level; child issues cover the individual pieces of work. Not every issue needs sub-issues — only when a single issue would otherwise try to cover too much. Prefer several focused issues under a parent over one large issue with a long description.

## Milestones

Milestones map to versions (e.g. `v0.2.0`). Check the active milestone before picking up work — it tells you what the current focus is. Issues not assigned to a milestone are backlog.

## Commits and versioning

Use conventional commits — release-please reads them to determine the next version:

| Prefix   | Version bump  |
| -------- | ------------- |
| `feat:`  | minor (0.x.0) |
| `fix:`   | patch (0.0.x) |
| `feat!:` | major (x.0.0) |
| `chore:` | none          |

Do not add `Co-Authored-By` trailers to commits.

## Branch and PR flow

Trunk-based: `main` is the single long-lived branch and is always releasable.

- Branch off `main` for every change; keep branches short-lived
- PRs merge into `main` via squash merge — the PR title becomes the commit, so it must follow conventional commits format
- Merging a feature PR does **not** cut a release. release-please batches merged changes into a pending version PR; a release ships only when that version PR is merged
- release-please maintains the version PR automatically; do not edit it manually

### Testing changes together before merge

To validate several in-flight branches as a combination without merging any of them to `main`, create a throwaway integration branch, merge the branches into it, and test there:

```sh
git switch -c integration main
git merge feature-a feature-b
# build and test the combination
git switch main && git branch -D integration
```

The integration branch is disposable — never open a PR from it and never promote it. Delete it once testing is done.

## Release pipeline

`conventional commit → main` → release-please maintains the version PR → merge the version PR → tag pushed → GoReleaser builds macOS binaries → Homebrew formula updated automatically

## Configuration and environment variables

The configuration and setup model is designed in #54 and its child issues; consult it
before adding any configuration surface. Two policies hold regardless of how that design
evolves:

- **Environment variables are for operational, cross-cutting concerns and for selecting
  file locations — never for product settings.** Operational means diagnostics and
  terminal/output behaviour (e.g. `DEMI_DEBUG`, `NO_COLOR`); path-selectors choose where
  demi reads and writes (e.g. `DEMI_CONFIG_DIR`, `DEMI_STATE_DIR`) and are resolved at
  startup, ahead of any configuration. What demi tracks, and how it behaves as a product,
  is configured through demi — not the environment. Do not add an environment variable
  for a product setting.
- **Read the environment once into typed, immutable values behind an injectable lookup —
  never scatter `os.Getenv`.** This keeps behaviour deterministic and tests hermetic: a
  test supplies its own values instead of depending on the ambient environment.
  `internal/config` and `internal/iostreams` hold the authoritative set of variables.

Keep the environment surface small — a new variable must fall into one of the two
categories above.

## Debug output

Debug output is a single env-gated diagnostic stream on stderr. `DEMI_DEBUG` set to any
non-empty value other than `0`, `false`, or `no` (case-insensitively) turns it on;
unset or empty leaves it off.

- **Every debug line goes through `f.IOStreams.Debugf`** — it gates on the variable and
  prefixes `DEBUG: `. Never hand-roll the gate, the prefix, or the stream: no direct
  `os.Stderr` writes, no `fmt.Print*` or bare `print`/`println`, no logging packages
  (`log`, `log/slog`, or any third-party logger). The linter enforces the mechanical
  subset — process streams, print builtins, ambient environment reads, and the stdlib
  logging packages fail `make lint`; hand-rolling the gate or prefix against
  `IOStreams.Err` passes lint and is still forbidden.
- **`DEMI_DEBUG` is read in exactly one place** — `IsDebug()` in `internal/iostreams`.
  Call sites consume the accessors and never read the environment.
- **Debug lines are additive diagnostics, not presentation.** They may appear on success
  and on failure, but they never replace, decorate, or reorder error output — errors
  still surface exactly once, through `printError` (see Error handling).
- **Call sites live at genuine decision points** — a resolved path, a selected provider,
  a branch taken. Add them alongside the code they describe, not speculatively. When
  debug is on, the stream's first line is the entrypoint's version-and-args line.
- **The line format is not a contract.** Pre-1.0 it may change without deprecation;
  nothing may parse it.

That is the entire mechanism: one method, one level, one stream, one gate. No verbosity
flags, no levels, no structured logging — extend it only when a concrete subsystem needs
more than an on/off stream, and design that extension against the consumers that exist
then.

## Error handling

Errors flow up; the entrypoint (`cmd/demi` and `internal/demi`) owns presentation and
exit codes. Code below the entrypoint constructs and returns errors — it never prints
them and never calls `os.Exit`.

### Constructing

- **Never discard a cause.** When an underlying call fails, carry it:
  `fmt.Errorf("determine home directory for config: %w", err)` — not a fresh
  `errors.New` that loses the reason.
- **Wrap with `%w` by default.** Use `%v` only to deliberately sever the chain — when
  the cause is an implementation detail no caller may ever match on. Severing is the
  exception; when unsure, wrap. A deliberate sever carries `//nolint:errorlint` with a
  brief reason, so every sever is visible and greppable.
- **Messages are lowercase, with no trailing punctuation and no `error:` prefix**; they
  state what could not be done. Control sentinels whose text never prints (`ErrSilent`,
  `ErrCancel`) are named after their identifier instead — leave them as they are.
- **Add a sentinel** (`var ErrFoo = errors.New(...)`) only when a caller branches on
  identity; **add an error type** only when a caller needs data or behaviour from it.
  Otherwise return plain wrapped errors. Sentinels and types the entrypoint or multiple
  packages branch on live in `internal/cmdutil`.

### Propagating

- **Wrap only where you add information the callee lacked** — the path, the resource,
  the intent. If a wrap would restate the callee's message, `return err` unchanged.
- **The failure phrase (`could not …`) is added once, by the layer nearest the user —
  normally the command; wraps below it name the operation** (`read settings file: %w`),
  so the final message reads as one sentence — never `could not X: could not Y`.
- **Classification sees the whole chain.** `errors.Is`/`errors.As` match through every
  `%w` link, so wrapping republishes everything inside: a returned chain containing
  `context.Canceled` exits 2 with nothing printed — the wrap's message is discarded.
  If a dependency leaks `context.Canceled` from a non-cancellation failure, sever that
  link with `%v`; never sever a genuine cancellation — that turns Ctrl-C into a
  printed exit-1 failure.

### Presentation and exit codes

- **A returned error surfaces exactly once**: `printError` in `internal/demi` prints it
  — code never both prints an error and returns it. It writes to stderr and appends
  usage for flag and unknown-command errors.
- **A command that has already shown its failure returns `cmdutil.ErrSilent`** —
  presenting the cause is not discarding it. Non-fatal diagnostics a command emits
  while succeeding are ordinary output, outside this contract.
- **A command that stops because its context was canceled returns `ctx.Err()`**, bare
  or wrapped — `mapError` recognises it chain-deep and stays silent. `cmdutil.ErrCancel`
  is for user aborts that carry no context error (e.g. a declined prompt).
- **`mapError` in `internal/demi` is the only place exit codes are decided**: nil → 0,
  failure → 1, user cancellation → 2. `context.DeadlineExceeded` is a failure, not a
  cancellation. Exit codes are a documented user contract — a new code means a new
  sentinel or type in `internal/cmdutil`, a `mapError` branch, and its docs.
- **Classify with `errors.Is`/`errors.As` against named sentinels and types — never by
  matching message text.** The entrypoint's unknown-command prefix check is the single
  deliberate classification exception; do not add another.

### What not to add

Plain stdlib errors are the entire mechanism. No error-handling packages, no error
codes, no stack traces, no structured logging — nothing of the kind until a concrete
consumer exists. The intact `%w` chain is the future diagnostic surface; preserving it
is what the rules above buy.

## Quick context commands

```sh
gh issue list --state open                          # what is open
gh api repos/{owner}/{repo}/milestones              # current version goals
gh pr list                                          # open PRs including the release-please PR
git log --oneline -10                               # recent commits
```

## Testing

Run the full suite before opening a PR:

```sh
make test   # go test -race ./...
make lint   # golangci-lint run (requires golangci-lint v2 locally: brew upgrade golangci-lint)
make check  # runs both; mirrors the required CI checks
```

**Fakes:** use `iostreams.Test()` to get pre-wired in/out/err streams for unit tests — do not reach for real file descriptors or `os.Stdout`.

**Environment pinning:** any test that exercises code reading the ambient environment must pin every variable that code reads with `t.Setenv`. Tests that supply values through an injected lookup do not need `t.Setenv` for those values — prefer injection where a seam exists.

**End-to-end (testscript):** the CLI harness lives in `internal/demi/main_test.go`; scripts go under `internal/demi/testdata/scripts/`. Each script is a self-contained `testscript` scenario. Add a new `.txtar` file there for any behaviour that is easier to verify at the binary level than at the unit level.
