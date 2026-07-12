# Demiurge — Agent Guide

## What this project is

A macOS CLI tool (`demi`) for developers who maintain many projects simultaneously. It provides terminal-native project awareness — not a project manager, no tasks or boards. Local-only, no accounts, no cloud. See `README.md` for full context.

## Before starting any work

Run this to get current state:

```sh
gh issue list --state open
gh issue list --milestone <current-milestone> --state open
gh milestone list
```

**Every piece of work must be linked to at least one open issue.** If there is no issue that covers what you are about to do, stop and ask the user whether to create one before proceeding. Do not open issues yourself without asking first.

## Issue workflow

- All work (features, fixes, chores that change behaviour) links to an issue
- Reference issues in commit messages: `feat: add git status view, closes #7`
- Multiple related commits can reference the same issue without closing it: `ref #7`
- Issues close automatically when a release ships if `closes #N` appears in a commit that lands on `main`

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

`conventional commit → main` → release-please maintains the version PR → merge the version PR → tag pushed → GoReleaser builds macOS binaries → Homebrew cask updated automatically

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

## Quick context commands

```sh
gh issue list --state open          # what is open
gh milestone list                   # current version goals
gh pr list                          # open PRs including the release-please PR
git log --oneline -10               # recent commits
```
