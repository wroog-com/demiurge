# Vision & Domain Model

This document is for people building `demi` (human or agent). It records *what
demi is* and the reasoning behind the domain model, so that design decisions are
made against a shared understanding rather than re-derived each time.

It deliberately stops short of prescribing the provider/dependency architecture.
Design decisions and their progress are tracked in the repository's issues, and
should stay there until we decide otherwise — this document gives them the missing
context (*the product they serve*), it does not resolve them. Architectural
decisions in particular are meant to emerge from the first concrete integrations,
not from this document.

### Tracking design work

Design discussions and decisions are recorded as issues, not in prose documents,
so that the history and current status stay searchable in one place. To make
design work discoverable, label such issues with **`design`**. Anyone looking for
the design history or open architectural questions can then filter on that label:

```sh
gh issue list --label design --state all
```

A design issue closes when its decision and rationale are recorded on the issue
itself. Keep the reasoning in the issue; keep only durable, settled product
context here.

## What demi is

`demi` is a **project lifecycle companion**. It sits around a developer's normal
workflow and assists at the boundaries between projects. The developer's actual
work — writing code, running git, using an editor — is unchanged. `demi` is
active at the transitions: conception, creation, provisioning, maintenance,
retirement, and memory.

The single idea that makes all of this cohere: **demi accumulates meta-knowledge
about all of a developer's projects over time**, and that accumulated context is
what makes every stage smarter. Without it, `demi` is just a wrapper around some
commands. With it, `demi` can answer "haven't I done this before?", remind you
after teardown, and reason across your whole body of work.

This has a consequence worth stating plainly: **demi is all-in at the edges,
hands-off in the middle.** Its value is proportional to how consistently the
lifecycle boundaries — creating, provisioning, retiring — are routed through it.
A project created outside `demi` cannot be tracked; setup performed outside `demi`
cannot be reversed by it; and every skipped step is a gap in its cross-project
memory. So `demi` only works well when it's used as the way you handle those
boundaries. Note the boundary of that claim: the *inner* work — writing code,
running git, using an editor — is deliberately **not** routed through `demi` and
stays exactly as it is. The commitment is to the transitions, not to the daily
work between them.

## What a "project" is

A project in `demi` is **a registry entry that demi maintains** — a tracked thing
with metadata — *not* a body of data that demi owns or stores on your behalf.

This is the most important distinction in the whole design. `demi` does not own
project *contents* — it does not house your code, and it is not a storage backend
for it:

- `demi` provides **awareness**, not custody. Your code already lives in the
  filesystem and in git. `demi` reads and acts on it; it does not house it.
- "Creating" a project therefore means **starting to track it** — adding it to
  `demi`'s list and capturing some metadata — with no external side effects until
  you explicitly ask for provisioning.

So the thing `demi` persists is small: the list of tracked projects, their
metadata, and a record of the setup `demi` performed for each. This is `demi`'s
own state. It lives locally (see `internal/config` — `StateDir()`), and it is
always present regardless of any integration.

## Load-bearing constraint: clean reversibility

The retirement promise — *remove the project from all places, leaving no trail* —
is the constraint that most shapes the architecture, and it is easy to overlook
until it's too late to add cheaply.

It only holds if `demi` records **provenance** for every external side effect it
performs: not merely "this project exists," but "*demi* created this hosted repo,
*demi* wrote this `.gitignore`, *demi* injected these files." Teardown then
reverses precisely what `demi` did and nothing it didn't.

Two things follow:

1. **This is a product principle, stated now.** Every future integration must be
   designed so its actions are *recordable and reversible*. An integration that
   can set something up but cannot cleanly undo it violates the core promise.
2. **The mechanism is not decided here.** *How* provenance is recorded (a ledger,
   event log, per-action metadata, something else) is an implementation decision
   that belongs with the provider/architecture design work tracked in the issues,
   driven by the first real integrations. This document commits to the *behavior*,
   not the *mechanism*.

Note also that this is **not** a storage-backend concern. The pluggable seam in
`demi` — if and when one is justified — is about *what external systems demi can
act on and later reverse* (code hosting, filesystem operations on the project, AI,
…), not about *where demi keeps its own list*. `demi`'s own state is small, local,
and singular. The earlier "storage provider" instinct was aimed one layer off:
the interesting variation is in **actions on the outside world**, not in the
registry's storage medium.

## Guardrails for future work

- Don't let `demi` drift into being a project manager (tasks, boards, scheduling).
  It assists a workflow; it is not the workflow.
- Don't build the provider abstraction speculatively. Build concretely, extract
  the interface when the second implementation actually exists — wait for concrete
  integrations to drive it.
- Preserve reversibility: no integration ships an action it cannot cleanly undo.
- Keep `demi`'s own state local, small, and singular. Cross-machine sync of the
  registry, if ever wanted, is a distinct concern — not a storage-provider one.
