# CLAUDE.md - bipline

Operating instructions for Claude Code working in this repo. Read this before planning or editing.

## What this is
bipline turns GitHub shipping events (a published release, a merged PR) into reviewed build-in-public drafts in Buffer. It runs as a Go CLI invoked by a GitHub Action. Full spec in `docs/bipline-prd.md`.

## Stack
- Go (CLI, Buffer client, draft generator)
- Buffer GraphQL API (sink), via genqlient
- Claude API (drafting)
- GitHub Actions (trigger)
- Module path: `github.com/davigiroux/bipline`

## Hard invariants (never violate)
- Drafts only. Nothing auto-publishes or auto-schedules. The human review gate in Buffer is the point.
- Single user. Personal Buffer API key only. No third-party OAuth, no multi-tenant code.
- No daemon, no hosting, no database. CI is the only runtime.
- Generated post text follows the voice rules in `voice.md`. No em dashes anywhere.

## How to work here
- Stay within the current phase (see below). Do not build ahead into later phases.
- Propose a plan before editing files. Wait for approval.
- Small, real commits. One PR per phase.
- A phase is done only when its acceptance criteria in `ROADMAP.md` pass. Run the tests.
- Record non-obvious decisions in `DECISION-LOG.md`.

## Davi owns these, do not write them without him
- The `.graphql` operation files. These are the API contract and Davi hand-authors them.
- The voice prompt inside the generator. It is his voice, he writes it.
- The idempotency guard logic. Propose options, but he decides and writes it.

Everything else (CLI wiring, test scaffolding, `genqlient.yaml`, Actions YAML, struct plumbing) is fair to draft.

## Current phase
**Phase 0 - Spike (manual, no code yet).** <!-- Davi: update this line as you advance -->

## Buffer API context
Endpoint: `https://api.buffer.com`. Auth: `Authorization: Bearer $BUFFER_API_KEY`.
Object model: Account -> Organizations -> Channels -> Posts. Ideas belong to an Organization.

<!-- FILL IN after the Phase 0 spike -->
- Organization ID: 6a073bf15d53897094b8a76b
- Channels used for posting (service -> channel ID):
  - 
    ```
    {"data":{"channels":[{"id":"6a073c45090476fb9922fd92","name":"devgiroux","service":"instagram"},{"id":"6a073cec090476fb99230076","name":"devgiroux","service":"twitter"},{"id":"6a077599090476fb99243bdc","name":"davi-alvarenga-028614119","service":"linkedin"}]}}
    ```
- Draft path chosen (createPost with saveToDraft=true, OR createIdea): We'll go with ideas for now.
- Exact required arguments for the chosen draft mutation (copy from the API Explorer):
  ```graphql
  TODO
  ```

## Claude API context
- Model: `sonnet` <!-- Davi: pick the model string -->
- Voice constraints live in `voice.md`. The generator references them, do not duplicate them in code.

## Repo layout (target)
```
cmd/bipline/main.go        CLI entrypoint
internal/eventsource/      parse GITHUB_EVENT_PATH -> Event
internal/generator/        Event + voice -> draft (Claude API)
internal/buffer/           genqlient client + thin wrapper
voice.md                   voice constraints, single source of truth
.github/workflows/bipline.yml
docs/bipline-prd.md        full spec
```

## Commit & PR conventions
- Imperative, scoped commit subjects (e.g. `buffer: handle MutationError branch`).
- One PR per phase, titled `Phase N: <name>`.
- Keep generated code (genqlient bindings) in its own commit, separate from hand-written code, so the history shows what is yours.
