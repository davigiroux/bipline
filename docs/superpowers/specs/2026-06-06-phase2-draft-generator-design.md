# Phase 2: Draft Generator — Design Spec

**Date:** 2026-06-06
**Status:** Approved
**Branch:** davigiroux/phase-2
**Scope:** Phase 2 from ROADMAP.md. Generator package only. No Buffer calls, no CLI wiring, no event parsing.

---

## Context

Phase 1 delivered a typed Buffer client. Phase 2 adds the LLM layer: given a normalized event and a voice file, produce a social post draft string. The generator is pure — no network calls to Buffer, no file I/O except reading `voice.md`. The CLI (Phase 4) wires everything together.

---

## Package Layout

Two packages touched in Phase 2:

```
internal/eventsource/
└── event.go              # Event struct — type definition only (Phase 3 adds parsing)
internal/generator/
├── generator.go          # Draft() function + Claude API call
└── generator_test.go     # integration tests (//go:build integration)
voice.md                  # user writes this; placeholder committed in Phase 2
Makefile                  # add test-integration target
```

---

## Event Struct (`internal/eventsource/event.go`)

Defined in Phase 2 as the shared type between generator (consumer) and eventsource (producer). Phase 3 adds parsing functions to the same package.

```go
// Event is a normalized GitHub shipping event.
type Event struct {
    Type  string // "release"
    Repo  string // "owner/repo"
    URL   string // canonical link to the release
    Title string // release name
    Body  string // release notes or description
    Tag   string // tag name, e.g. "v1.2.0"
}
```

Fields are minimal — only what the generator needs to write a post. The `Type` field is a string for now; Phase 3 will determine if it needs to become a typed constant.

---

## Generator Function (`internal/generator/generator.go`)

```go
// Draft generates a social media post draft for the given event.
// voicePath is the path to the user's voice.md file.
// ANTHROPIC_API_KEY must be set in the environment.
func Draft(ctx context.Context, event eventsource.Event, voicePath string) (string, error)
```

**Internals:**

1. Read `voice.md` from `voicePath` — full file contents become the Claude system prompt verbatim. No parsing, no structured format.
2. Build user message from event fields — title, body, repo, URL, tag composed into a clear prompt.
3. Call Claude via `github.com/anthropics/anthropic-sdk-go`, model `claude-sonnet-4-6`, max tokens 500.
4. Return the text content of the first response content block.

The Anthropic client is constructed inside `Draft` from `ANTHROPIC_API_KEY` env var. No config structs, no dependency injection — function is the public API, Claude call is the implementation detail.

**Error handling:** wrap all errors with context (`fmt.Errorf("generator: read voice: %w", err)` etc.).

---

## voice.md Contract

The generator reads `voice.md` at call time and passes its contents as the system prompt to Claude. No processing — what's in the file is what Claude receives.

`voice.md` is a **manual handoff**: each user writes their own. A placeholder is committed with:
- Annotated sections showing what to fill in
- Examples of the kind of constraint that belongs in each section (persona, style rules, format constraints, what to include)
- A note that this file is the single source of truth for post voice — do not duplicate in code

The placeholder is instructive for any user, not opinionated about any specific voice.

---

## Testing

**Integration tests** (`//go:build integration`) in `generator_test.go`. Require `ANTHROPIC_API_KEY` in env (`.env` via Makefile).

Two fixture `Event` values representing example past releases. Marked clearly as user-replaceable examples:

```go
// Replace these with your own past events before running integration tests.
var fixtureRelease = eventsource.Event{
    Type:  "release",
    Repo:  "owner/your-repo",
    URL:   "https://github.com/owner/your-repo/releases/tag/v0.1.0",
    Title: "v0.1.0 - Initial release",
    Body:  "Short description of what shipped.",
    Tag:   "v0.1.0",
}
```

**Assertions on each generated draft:**
- No em dash (`—`) anywhere
- Contains the event URL
- Length ≤ 500 characters

These are invariants, not exact-match assertions — Claude's output is non-deterministic.

**Makefile target added:**
```makefile
test-integration:
	go test -tags integration -v ./internal/generator/
```

---

## Multi-User Design

The generator is generic by design:
- `voicePath` is a parameter — each user points to their own voice file
- `ANTHROPIC_API_KEY` comes from env — standard secret injection pattern
- Event struct is defined by the caller — no project-specific defaults

Fixture events in tests are clearly marked as examples. The `voice.md` placeholder documents the structure without prescribing a specific voice.

Note for Phase 4: the Buffer org ID and channel IDs currently hardcoded in `cmd/probe/main.go` must become CLI flags or a config file before Phase 4 is considered done. This is the main remaining multi-user concern.

---

## Ownership Summary

| File | Owner |
|------|-------|
| `internal/eventsource/event.go` | Claude scaffolds |
| `internal/generator/generator.go` | Claude writes |
| `internal/generator/generator_test.go` | Claude writes (fixtures are user-replaceable) |
| `voice.md` | **User writes** (Claude commits instructive placeholder) |
| `Makefile` (test-integration target) | Claude adds |

---

## Acceptance Criteria (from ROADMAP.md)

- [ ] `generator.Draft(event, voice)` returns a draft string
- [ ] Golden tests on 2-3 real past events pass voice invariants (no em dash, contains URL, ≤ 500 chars)
- [ ] No Buffer network calls in this package

---

## Verification

1. `go build ./...` — clean
2. Davi writes `voice.md`
3. Update fixture events in `generator_test.go` with real past events
4. `make test-integration` — both tests pass, output contains repo URL, no em dashes
