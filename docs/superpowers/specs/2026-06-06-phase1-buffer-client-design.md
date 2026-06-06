# Phase 1: Buffer Client — Design Spec

**Date:** 2026-06-06  
**Status:** Approved  
**Scope:** Phase 1 from ROADMAP.md. Buffer package only. No generator, no CLI, no GitHub Action.

---

## Context

Phase 0 (manual spike) is complete. Org ID, channel IDs, and draft path (createIdea) are confirmed in CLAUDE.md. Phase 1 wires up a typed Go client backed by genqlient, proves the full round-trip (resolve org → find channel → create idea), and handles both success and error return branches from the mutation.

---

## Repository Setup

- Create public GitHub repo: `davigiroux/bipline`
- Initialize Go module: `github.com/davigiroux/bipline`
- One PR for Phase 1, titled `Phase 1: Buffer client`
- Generated genqlient bindings committed separately from hand-written code

---

## Directory Layout (Phase 1 only)

```
internal/buffer/
├── genqlient.yaml       # codegen config (scaffolded by Claude)
├── schema.graphql       # fetched via introspection (see below)
├── operations.graphql   # Davi writes; placeholder committed first
├── generated.go         # output of go generate; committed alone
└── client.go            # thin wrapper (Claude writes)
internal/buffer/README.md
cmd/probe/main.go        # integration smoke test (Claude writes)
tools/fetch-schema.sh    # one-shot introspection script (Claude writes)
Makefile                 # targets: schema, generate, probe
```

---

## Schema Fetch

`tools/fetch-schema.sh` runs a GraphQL introspection query via curl against `https://api.buffer.com`, using `$BUFFER_API_KEY`. Output goes to `internal/buffer/schema.graphql`. Davi runs this once; the file is committed.

Makefile target: `make schema`

---

## genqlient Config

`genqlient.yaml`:
```yaml
schema: schema.graphql
operations:
  - operations.graphql
generated: generated.go
package: buffer
```

---

## operations.graphql (Davi's file)

Claude commits a placeholder with the expected operations documented inline. Davi replaces with correct GraphQL syntax. Expected operations based on PRD:

- `GetChannels(organizationId: String!)` → channels list (id, name, service)
- `CreateIdea(organizationId: String!, title: String!, text: String!)` → success/error union

The placeholder will have these as comments with the expected field shapes so Davi can fill in the syntax.

---

## BufferClient (client.go)

```go
type BufferClient struct {
    apiKey string
    orgID  string
    client *http.Client
}

type Channel struct {
    ID      string
    Name    string
    Service string
}

func New(apiKey, orgID string) *BufferClient

// FindChannel returns the channel matching the given service name.
// Returns error if not found.
func (c *BufferClient) FindChannel(ctx context.Context, service string) (*Channel, error)

// CreateIdea creates a draft idea in the org's Ideas board.
// Handles both success and MutationError branches from the API response.
func (c *BufferClient) CreateIdea(ctx context.Context, title, text string) error
```

`CreateIdea` inspects the mutation response union and returns a typed `MutationError` value if the API returns the error branch. Davi writes the operation; Claude writes the Go-side branch handling once generated types are available.

---

## Smoke Test (cmd/probe/main.go)

Reads `BUFFER_API_KEY` from env. Instantiates `BufferClient` with hardcoded org ID from CLAUDE.md. Calls `FindChannel("twitter")`, then `CreateIdea("probe test", "testing the bipline buffer client")`. Prints result. Not a unit test — this is the acceptance-criteria validator for Phase 1.

---

## Unit Tests

One file: `internal/buffer/client_test.go`. Table-driven test for `FindChannel`: stubs the HTTP transport to return a fixture channels response, asserts correct channel is returned for a given service name, and asserts error on unknown service.

---

## Ownership Summary

| File | Owner |
|------|-------|
| `genqlient.yaml` | Claude scaffolds |
| `schema.graphql` | fetched (Davi runs `make schema`) |
| `operations.graphql` | **Davi writes** (Claude commits placeholder) |
| `generated.go` | genqlient (Davi runs `make generate`) |
| `client.go` | Claude writes |
| `cmd/probe/main.go` | Claude writes |
| `tools/fetch-schema.sh` | Claude writes |
| `Makefile` | Claude writes |
| `internal/buffer/README.md` | Claude writes |

---

## Acceptance Criteria (from ROADMAP.md)

- [ ] genqlient generates bindings from the live schema
- [ ] Small program resolves org, finds channel by service, creates a draft
- [ ] Both mutation result branches (success + error) handled in client.go
- [ ] Package README committed

---

## Verification

1. `make schema` — succeeds, `schema.graphql` populated
2. Davi writes `operations.graphql`
3. `make generate` — `generated.go` produced, `go build ./...` clean
4. `BUFFER_API_KEY=... go run ./cmd/probe/` — creates an idea in Buffer
5. `go test ./internal/buffer/` — passes
