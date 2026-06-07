# Phase 4: CLI + GitHub Action — Design Spec

**Date:** 2026-06-06
**Status:** Approved
**Branch:** davigiroux/phase-4
**Scope:** Phase 4 from ROADMAP.md. Wires eventsource → generator → buffer into a CLI. Commits a GitHub Actions workflow to SafeNudge. Idempotency via Actions cache.

---

## Context

Phases 1-3 built the three isolated packages. Phase 4 glues them together into a real CLI (`cmd/bipline/main.go`) and deploys it to SafeNudge via a GitHub Actions workflow. A release on SafeNudge triggers the workflow, bipline drafts a post, and the draft lands in Buffer for review. A re-run of the same release produces zero new drafts.

---

## CLI (`cmd/bipline/main.go`)

### Config — all via env vars

| Env var | Source | Purpose |
|---|---|---|
| `GITHUB_EVENT_NAME` | GitHub Actions (automatic) | Event type: "release" or "pull_request" |
| `GITHUB_EVENT_PATH` | GitHub Actions (automatic) | Path to the raw event payload JSON |
| `BUFFER_API_KEY` | Secret | Buffer auth |
| `ANTHROPIC_API_KEY` | Secret | Claude auth |
| `BUFFER_ORG_ID` | Var | Buffer organization ID |
| `BUFFER_SERVICE` | Var | Target channel service, e.g. "twitter" |
| `BIPLINE_VOICE_PATH` | Var | Path to voice.md |

All six non-automatic vars are required. The CLI validates them upfront and exits 1 with a clear message listing any missing vars before making any API calls.

### Pipeline

```
os.ReadFile(GITHUB_EVENT_PATH)
  → eventsource.Parse(GITHUB_EVENT_NAME, payload)
  → generator.Draft(ctx, event, BIPLINE_VOICE_PATH)
  → buffer.New(BUFFER_API_KEY, BUFFER_ORG_ID).CreateIdea(ctx, event.Title, draft)
```

### Exit codes

- `0` — success (idea created)
- `0` — unhandled event type or unmerged PR (`Parse` returns unsupported/unmerged error) — Action stays green on irrelevant events
- `1` — missing env vars, file read error, API error

### Unit-testable helper

`requireEnv(vars ...string) (map[string]string, error)` — checks each name is non-empty, returns all values or a single error listing every missing var. Tested in `cmd/bipline/main_test.go`.

---

## GitHub Actions Workflow (committed to SafeNudge)

File: `.github/workflows/bipline.yml` in the `davigiroux/safenudge` repository.

```yaml
name: bipline

on:
  release:
    types: [published]

jobs:
  draft:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/checkout@v4
        with:
          repository: davigiroux/bipline
          path: bipline

      - uses: actions/setup-go@v5
        with:
          go-version-file: bipline/go.mod

      - name: Restore idempotency cache
        id: cache
        uses: actions/cache@v4
        with:
          path: .bipline-cache
          key: bipline-${{ github.event.release.tag_name }}

      - name: Run bipline
        if: steps.cache.outputs.cache-hit != 'true'
        working-directory: bipline
        env:
          BUFFER_API_KEY: ${{ secrets.BUFFER_API_KEY }}
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          BUFFER_ORG_ID: ${{ vars.BUFFER_ORG_ID }}
          BUFFER_SERVICE: ${{ vars.BUFFER_SERVICE }}
          BIPLINE_VOICE_PATH: ./voice.md
          GITHUB_EVENT_PATH: ${{ github.event_path }}
          GITHUB_EVENT_NAME: ${{ github.event_name }}
        run: go run ./cmd/bipline/

      - name: Mark processed
        if: steps.cache.outputs.cache-hit != 'true'
        run: mkdir -p .bipline-cache && echo "ok" > .bipline-cache/${{ github.event.release.tag_name }}
```

**How idempotency works:** `actions/cache` restores `.bipline-cache` keyed on the release tag. On a re-run, `cache-hit == 'true'` and both "Run bipline" and "Mark processed" are skipped via `if:` conditions. First run: cache miss → bipline runs → marker written → cache saved. Second run: cache hit → both steps skipped → zero new drafts.

**voice.md:** Read from the bipline repo checkout (`./voice.md` with `working-directory: bipline`). SafeNudge carries no voice file — single source of truth stays in the bipline repo.

---

## Testing

| Test | How |
|---|---|
| `requireEnv` unit test | `go test ./cmd/bipline/` — no network, no API |
| CLI compiles | `go build ./cmd/bipline/` |
| Live acceptance | Cut a test release on SafeNudge; confirm one idea in Buffer. Re-run workflow; confirm zero new ideas. |

No unit tests for `main()` itself — it's wiring code and mocking all three packages would test nothing real.

---

## Ownership

| File | Owner |
|---|---|
| `cmd/bipline/main.go` | Claude writes |
| `cmd/bipline/main_test.go` | Claude writes |
| `.github/workflows/bipline.yml` (in SafeNudge) | Claude drafts — Davi commits to SafeNudge |
| Idempotency logic | Davi owns per CLAUDE.md — implemented as Actions cache per this design |

---

## Acceptance Criteria (from ROADMAP.md)

- [ ] CLI reads `GITHUB_EVENT_PATH`, drafts, creates a Buffer draft
- [ ] Workflow committed in SafeNudge
- [ ] A test release produces exactly one draft
- [ ] A re-run produces zero new drafts (idempotency works)
