# Demiurge

`demi` is a terminal-native **project lifecycle companion** for developers who
maintain many projects at once.

It is not a project manager — there are no tasks, boards, or tickets. It does not
replace your editor, your shell, or git. Your day-to-day work stays exactly as it
is. `demi` sits _around_ that work, at the moments between projects, and makes the
edges smarter: deciding what to build, setting it up, keeping it tidy, tearing it
down cleanly, and remembering what you've done.

Local-only. No accounts, no cloud, no sync required.

## Install

`demi` currently runs on macOS — both Apple Silicon and Intel.

Install with Homebrew:

```sh
brew install wroog-com/tap/demi
```

Or build from source with Go:

```sh
go install github.com/wroog-com/demiurge/cmd/demi@latest
```

## Why it exists

If you keep a lot of projects going, the hard parts usually aren't in the code —
they're around it:

- Starting a new project means the same setup ritual every time.
- Old experiments pile up, half-configured, and you're never sure what's safe to
  delete.
- When you do delete one, you're never sure you got _everything_ — the repo, the
  local files, the bits of config scattered around.
- You have an idea and can't remember whether you already tried something like it,
  or you do remember you did but have no idea where to find it.

`demi` exists to hold that context for you. Because it knows about all your
projects, it can help at each stage instead of leaving you to remember it all.

## The lifecycle

`demi` thinks about a project as something that moves through stages. Between the
stages, you're just working — editor, terminal, git, as usual. At the stages,
`demi` helps.

1. **Conceive** — Before you start, ask `demi` whether an idea is actually new.
   Because it knows your whole history, it can tell you "you already have two
   projects that look like this" and give you the context to decide.
2. **Create** — Register the project so `demi` starts tracking it. This is
   lightweight: it adds the project to the list it maintains and asks a few
   questions to capture what it is. Nothing external happens yet.
3. **Provision** — When you're ready, have `demi` wire up the real-world pieces:
   create the hosted repository, apply a `.gitignore` from your rules, drop in the
   files and tooling you want every project to start with.
4. **Maintain** — Reset a project that's gone sideways, spin up a clean "v2", or
   re-apply your standard setup as your conventions evolve.
5. **Retire** — When you're done, `demi` cleans up after itself. Because it
   tracked what it set up, it can undo exactly that — leaving no trail — without
   guessing.
6. **Remember** — On the way out, `demi` can keep what it learned about the
   project, so the next time you _conceive_ something similar it can remind you.

These stages are a way of thinking about the lifecycle, not a fixed checklist.
They are non-exhaustive and will adapt as `demi` grows — expect them to be
refined, split, or added to over time.

The AI-assisted parts of this (the "is this new?" reasoning, opinions about your
projects) are a direction the tool is built toward, not all present yet.

## Principles

- **Companion, not manager.** `demi` assists a normal workflow; it never becomes
  the workflow.
- **`demi` owns only what `demi` created.** It tracks your projects and the setup
  it performed on your behalf. It does not take ownership of your code — the
  filesystem and version control (and many other systems) already do that.
- **Reversible by design.** Anything `demi` sets up, `demi` can cleanly take back
  down. That promise is what makes it safe to let it do setup work at all.
- **Local-first.** Everything works on one machine with no external services. Any
  integration (code hosting, AI, and so on) is something you opt into.
- **Accumulated context is the point.** The value grows as `demi` learns more
  about your projects over time.

## Stability

`demi` is pre-1.0 and under active development. Anything you might depend on —
commands, flags, environment variables, and stored state — may change without a
deprecation period until 1.0.

## Reference

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error — message written to stderr |
| 2 | Cancelled — Ctrl-C or SIGTERM; nothing printed |

Commands that cooperate with cancellation exit 2 on the first Ctrl-C or SIGTERM,
rather than dying by signal (130 or 143). A second signal after the first has
been handled kills the process; the shell sees 130 (SIGINT) or 143 (SIGTERM). A
gracefully-cancelled command that exits 2 does not abort an enclosing bash `for`
loop the way signal-death would — the escape applies only to a process still
running.

### Environment variables

| Variable | Effect |
|----------|--------|
| `DEMI_DEBUG` | Set to any non-empty value (not `0`, `false`, or `no`) to write diagnostic output to stderr |

## For contributors

- [`docs/vision.md`](docs/vision.md) — the product model and the reasoning behind
  it, written for people building `demi`.
- [`AGENTS.md`](AGENTS.md) — the working agreement: issue workflow, branching,
  releases. (`CLAUDE.md` is a thin `@AGENTS.md` import so tools expecting that
  filename still load it.)

Current focus and progress live in the repository's issues and milestones, not in
this file.
