# Phase 3: Event Source — Design Spec

**Date:** 2026-06-06
**Status:** Approved
**Branch:** davigiroux/phase-3
**Scope:** Phase 3 from ROADMAP.md. Parser only — no CLI wiring, no trigger logic, no label gate.

---

## Context

Phase 2 delivered the draft generator. `eventsource.Event` already exists as a struct in `internal/eventsource/event.go`. Phase 3 adds `Parse()` — a function that turns a raw GitHub JSON payload into that struct. Phase 4 (CLI) will call `Parse` after reading `GITHUB_EVENT_PATH` from the environment.

---

## Package Layout

Phase 3 adds to the existing `internal/eventsource/` package. No new packages.

```
internal/eventsource/
├── event.go          # existing — Event struct (unchanged)
├── parse.go          # new — Parse() function
└── testdata/
    ├── release.json  # GitHub release.published fixture
    └── pr.json       # GitHub pull_request (merged) fixture
```

---

## Parse Function

```go
// Parse parses a GitHub event payload into a normalized Event.
// eventName is the GitHub event name ("release" or "pull_request").
// payload is the raw JSON bytes from GITHUB_EVENT_PATH.
// Returns an error for unsupported event types or invalid payloads.
// For pull_request events, returns an error if the PR is not merged.
func Parse(eventName string, payload []byte) (Event, error)
```

The function takes bytes — no filesystem access. The CLI (Phase 4) reads `GITHUB_EVENT_PATH` and passes the bytes here. This keeps the parser pure and trivially testable.

---

## Field Mapping

| Event field | `release` payload source | `pull_request` payload source |
|-------------|--------------------------|-------------------------------|
| `Type`      | `"release"` (literal)    | `"pr"` (literal)              |
| `Repo`      | `repository.full_name`   | `repository.full_name`        |
| `URL`       | `release.html_url`       | `pull_request.html_url`       |
| `Title`     | `release.name`           | `pull_request.title`          |
| `Body`      | `release.body`           | `pull_request.body`           |
| `Tag`       | `release.tag_name`       | `""` (empty — PRs have no tag)|

---

## Error Cases

- **Unsupported event name:** returns `fmt.Errorf("eventsource: unsupported event %q", eventName)`. Covers anything other than `"release"` and `"pull_request"`.
- **Invalid JSON:** returns the `json.Unmarshal` error wrapped with context.
- **Non-merged PR:** `pull_request.merged == false` returns `fmt.Errorf("eventsource: pull_request is not merged")`. A closed-but-not-merged PR produces no post.

---

## Test Fixtures

Minimal JSON in `internal/eventsource/testdata/` — only fields bipline reads, not the full GitHub payload.

**`testdata/release.json`** (release.published):
```json
{
  "action": "published",
  "release": {
    "tag_name": "v1.2.0",
    "name": "v1.2.0 - Add batching",
    "body": "Groups notifications by sender. Cuts interruptions by ~60%.",
    "html_url": "https://github.com/devgiroux/safenudge/releases/tag/v1.2.0"
  },
  "repository": {
    "full_name": "devgiroux/safenudge"
  }
}
```

**`testdata/pr.json`** (pull_request, merged):
```json
{
  "action": "closed",
  "pull_request": {
    "merged": true,
    "title": "Add quiet hours setting",
    "body": "Lets users configure a window where notifications are suppressed.",
    "html_url": "https://github.com/devgiroux/safenudge/pull/42"
  },
  "repository": {
    "full_name": "devgiroux/safenudge"
  }
}
```

---

## Tests

One file: `internal/eventsource/parse_test.go`. Table-driven. No network, no API calls.

Test cases:
- `release.json` → correct Event fields
- `pr.json` → correct Event fields, Tag is empty string
- `pull_request` with `merged: false` → error
- Unsupported event name → error
- Invalid JSON → error

---

## Acceptance Criteria (from ROADMAP.md)

- [ ] Parses a `release.published` fixture into a normalized Event
- [ ] Parses a merged-PR fixture into a normalized Event
- [ ] Unit tested against fixture JSON
