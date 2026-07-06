c# Demiurge — Agent Guide

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

- Work happens on `dev` or feature branches
- PRs merge into `main` via squash merge — PR title becomes the commit, so it must follow conventional commits format
- release-please opens version PRs automatically on `main` pushes; do not modify those PRs manually

## Release pipeline

`conventional commit → main` → release-please opens version PR → merge PR → tag pushed → GoReleaser builds macOS binaries → Homebrew cask updated automatically

## Quick context commands

```sh
gh issue list --state open          # what is open
gh milestone list                   # current version goals
gh pr list                          # open PRs including release-please PRs
git log --oneline -10               # recent commits
```
