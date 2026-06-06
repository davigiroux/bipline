# Phase 1: Buffer Client Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a typed Go Buffer GraphQL client backed by genqlient that can create an Idea draft, satisfying Phase 1 acceptance criteria.

**Architecture:** genqlient codegen produces typed Go functions from `schema.graphql` + `operations.graphql`. A thin `BufferClient` wrapper adds auth and exposes clean `FindChannel`/`CreateIdea` methods. `cmd/probe` exercises the full round-trip as an integration smoke test.

**Tech Stack:** Go 1.22+, github.com/Khan/genqlient, Buffer GraphQL API at https://api.buffer.com

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `go.mod` | Create | Module declaration |
| `go.sum` | Create | Dependency lock |
| `.gitignore` | Create | Ignore build artifacts |
| `Makefile` | Create | schema / generate / test / probe targets |
| `tools/fetch-schema.sh` | Create | Introspection query → schema.graphql |
| `internal/buffer/generate.go` | Create | go:generate directive |
| `internal/buffer/genqlient.yaml` | Create | Codegen config |
| `internal/buffer/schema.graphql` | Fetched by Davi | Live Buffer API schema |
| `internal/buffer/operations.graphql` | Placeholder → Davi fills | GetChannels + CreateIdea ops |
| `internal/buffer/generated.go` | go generate output | Typed genqlient bindings |
| `internal/buffer/client.go` | Create | BufferClient + Channel + authTransport |
| `internal/buffer/client_test.go` | Create | FindChannel unit tests |
| `internal/buffer/README.md` | Create | Package docs |
| `cmd/probe/main.go` | Create | Integration smoke test |

---

### Task 1: Initialize GitHub repo and push planning docs

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Initialize git repo and add .gitignore**

In `/Users/davigiroux/projects/bipline`, run:

```bash
git init
```

Create `.gitignore`:

```
# binaries
/bipline
/probe

# test binaries
*.test

# env files
.env
.env.local
```

- [ ] **Step 2: Stage and commit existing planning docs**

```bash
git add CLAUDE.md DECISION-LOG.md ROADMAP.md docs/ .gitignore
git commit -m "initial: add planning docs and spec"
```

- [ ] **Step 3: Create public GitHub repo and push**

```bash
gh repo create davigiroux/bipline --public --source=. --remote=origin --push
```

Expected output: URL of the new repo printed, branch pushed.

---

### Task 2: Go module and genqlient dependency

**Files:**
- Create: `go.mod`, `go.sum`, `tools/tools.go`

- [ ] **Step 1: Initialize Go module**

```bash
go mod init github.com/davigiroux/bipline
```

Expected: `go.mod` created with `module github.com/davigiroux/bipline` and `go 1.22` (or current version).

- [ ] **Step 2: Add genqlient as dependency**

```bash
go get github.com/Khan/genqlient
```

Expected: `go.mod` and `go.sum` updated.

- [ ] **Step 3: Create tools.go to pin the generate tool**

Create `tools/tools.go`:

```go
//go:build tools

package tools

import _ "github.com/Khan/genqlient"
```

- [ ] **Step 4: Tidy and verify**

```bash
go mod tidy
go build ./...
```

Expected: no errors (nothing to build yet, but module is valid).

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum tools/
git commit -m "build: init Go module and add genqlient dependency"
```

---

### Task 3: Schema fetch script and Makefile

**Files:**
- Create: `tools/fetch-schema.sh`, `Makefile`

- [ ] **Step 1: Write fetch-schema.sh**

Create `tools/fetch-schema.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
: "${BUFFER_API_KEY:?BUFFER_API_KEY must be set}"

QUERY='{"query":"query IntrospectionQuery { __schema { queryType { name } mutationType { name } subscriptionType { name } types { ...FullType } directives { name description locations args { ...InputValue } } } } fragment FullType on __Type { kind name description fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason } inputFields { ...InputValue } interfaces { ...TypeRef } enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason } possibleTypes { ...TypeRef } } fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue } fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } } }"}'

mkdir -p internal/buffer

curl -sf -X POST https://api.buffer.com \
  -H "Authorization: Bearer $BUFFER_API_KEY" \
  -H "Content-Type: application/json" \
  -d "$QUERY" \
  | jq '.data' \
  > internal/buffer/schema.graphql

echo "schema.graphql written ($(wc -l < internal/buffer/schema.graphql) lines)"
```

- [ ] **Step 2: Write Makefile**

Create `Makefile`:

```makefile
.PHONY: schema generate test probe

schema:
	bash tools/fetch-schema.sh

generate:
	go generate ./internal/buffer/

test:
	go test ./...

probe:
	go run ./cmd/probe/
```

- [ ] **Step 3: Make script executable and commit**

```bash
chmod +x tools/fetch-schema.sh
git add tools/fetch-schema.sh Makefile
git commit -m "build: add schema fetch script and Makefile"
```

---

### ⚠️ MANUAL STEP: Fetch the live schema

> **Davi must run this.** `make schema` requires `BUFFER_API_KEY`.

```bash
BUFFER_API_KEY=<your-key> make schema
```

Expected: `internal/buffer/schema.graphql` populated with the Buffer GraphQL schema in introspection JSON format. The file should be non-empty (at minimum several hundred lines).

Once fetched:
```bash
git add internal/buffer/schema.graphql
git commit -m "buffer: add introspected Buffer GraphQL schema"
```

---

### Task 4: genqlient config and operations placeholder

**Files:**
- Create: `internal/buffer/generate.go`, `internal/buffer/genqlient.yaml`, `internal/buffer/operations.graphql`

> This task requires `internal/buffer/schema.graphql` to exist (from the manual step above).

- [ ] **Step 1: Create the go:generate directive file**

Create `internal/buffer/generate.go`:

```go
package buffer

//go:generate go run github.com/Khan/genqlient genqlient.yaml
```

- [ ] **Step 2: Create genqlient.yaml**

Create `internal/buffer/genqlient.yaml`:

```yaml
schema: schema.graphql
operations:
  - operations.graphql
generated: generated.go
package: buffer
```

- [ ] **Step 3: Commit the codegen config**

```bash
git add internal/buffer/generate.go internal/buffer/genqlient.yaml
git commit -m "buffer: add genqlient codegen config"
```

- [ ] **Step 4: Create operations placeholder**

Create `internal/buffer/operations.graphql`:

```graphql
# OWNER: Davi — replace the commented-out operations below with real GraphQL.
# Run `make generate` after editing this file.
#
# Expected operations based on docs/bipline-prd.md:
#
# query GetChannels($organizationId: String!) {
#   channels(input: { organizationId: $organizationId }) {
#     id
#     name
#     service
#   }
# }
#
# mutation CreateIdea($organizationId: String!, $title: String!, $text: String!) {
#   createIdea(input: {
#     organizationId: $organizationId
#     content: { title: $title, text: $text }
#   }) {
#     # Add fields here. Check schema.graphql for what createIdea returns.
#     # If the return type is a union, add __typename and fields for each variant.
#     # Example union pattern:
#     #   ... on IdeaActionSuccess { idea { id } }
#     #   ... on MutationError { message type }
#   }
# }
```

- [ ] **Step 5: Commit placeholder**

```bash
git add internal/buffer/operations.graphql
git commit -m "buffer: add operations.graphql placeholder (Davi to fill)"
```

---

### ⚠️ MANUAL STEP: Write operations.graphql

> **Davi must write this.** The `.graphql` operation files are hand-authored by Davi per CLAUDE.md.

1. Open `internal/buffer/operations.graphql`
2. Replace the commented-out stubs with real GraphQL operations (uncomment and correct field names against `schema.graphql`)
3. Ensure `GetChannels` and `CreateIdea` (or `CreateIdea`) match actual schema field names exactly
4. Check `schema.graphql` for the exact return type of `createIdea` — if it's a union, include `__typename` and all inline fragment fields

Once written:
```bash
git add internal/buffer/operations.graphql
git commit -m "buffer: write GetChannels and CreateIdea operations"
```

---

### Task 5: Run codegen and commit generated bindings

**Files:**
- Create: `internal/buffer/generated.go` (output of go generate)

> This task requires `operations.graphql` to be filled in by Davi.

- [ ] **Step 1: Run codegen**

```bash
make generate
```

Expected: `internal/buffer/generated.go` created. If it fails with a schema error, check that `operations.graphql` field names match `schema.graphql`.

- [ ] **Step 2: Verify the module builds**

```bash
go build ./...
```

Expected: clean build (nothing produces an executable yet, but generated code compiles).

- [ ] **Step 3: Note the generated types for use in Tasks 6 and 7**

Inspect `internal/buffer/generated.go`. Find and record:
- The name of the channels field in `GetChannelsResponse` (likely `Channels`)
- The type of each channel element (likely `GetChannelsChannels` with fields `Id`, `Name`, `Service`)
- The return type of `CreateIdea` — check if the `CreateIdea` field on `CreateIdeaResponse` is an interface (union) or a plain struct

- [ ] **Step 4: Commit generated code alone**

```bash
git add internal/buffer/generated.go
git commit -m "buffer: add genqlient-generated bindings [generated]"
```

---

### Task 6: BufferClient struct, FindChannel, and unit tests

**Files:**
- Create: `internal/buffer/client.go`, `internal/buffer/client_test.go`

- [ ] **Step 1: Write the failing test first**

Create `internal/buffer/client_test.go`:

```go
package buffer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
)

// fakeGQL implements graphql.Client by JSON-marshaling a fixed response into the
// resp.Data pointer that genqlient sets before calling MakeRequest.
type fakeGQL struct {
	data interface{}
	err  error
}

func (f *fakeGQL) MakeRequest(_ context.Context, _ *graphql.Request, resp *graphql.Response) error {
	if f.err != nil {
		return f.err
	}
	b, _ := json.Marshal(f.data)
	return json.Unmarshal(b, resp.Data)
}

func TestFindChannel(t *testing.T) {
	fixture := map[string]interface{}{
		"channels": []map[string]interface{}{
			{"id": "6a073cec090476fb99230076", "name": "devgiroux", "service": "twitter"},
			{"id": "6a077599090476fb99243bdc", "name": "davi-alvarenga-028614119", "service": "linkedin"},
		},
	}

	tests := []struct {
		service string
		wantID  string
		wantErr bool
	}{
		{service: "twitter", wantID: "6a073cec090476fb99230076"},
		{service: "linkedin", wantID: "6a077599090476fb99243bdc"},
		{service: "tiktok", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
			ch, err := c.FindChannel(context.Background(), tt.service)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ch.ID != tt.wantID {
				t.Errorf("got ID %q, want %q", ch.ID, tt.wantID)
			}
			if ch.Service != tt.service {
				t.Errorf("got Service %q, want %q", ch.Service, tt.service)
			}
		})
	}
}
```

- [ ] **Step 2: Run the test — expect compile error (BufferClient undefined)**

```bash
go test ./internal/buffer/ -run TestFindChannel -v
```

Expected output: compile error mentioning `BufferClient` and `FindChannel` undefined.

- [ ] **Step 3: Write client.go**

Create `internal/buffer/client.go`:

```go
package buffer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const endpoint = "https://api.buffer.com"

// Channel is a connected social media account.
type Channel struct {
	ID      string
	Name    string
	Service string
}

type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.base.RoundTrip(r)
}

// BufferClient is a thin wrapper over the genqlient-generated Buffer GraphQL client.
type BufferClient struct {
	orgID string
	gql   graphql.Client
}

// New creates a BufferClient authenticated with apiKey, scoped to orgID.
func New(apiKey, orgID string) *BufferClient {
	httpClient := &http.Client{
		Transport: &authTransport{apiKey: apiKey, base: http.DefaultTransport},
	}
	return &BufferClient{
		orgID: orgID,
		gql:   graphql.NewClient(endpoint, httpClient),
	}
}

// FindChannel returns the channel matching service (e.g. "twitter", "linkedin").
// Returns error if no channel is found for that service.
func (c *BufferClient) FindChannel(ctx context.Context, service string) (*Channel, error) {
	resp, err := GetChannels(ctx, c.gql, c.orgID)
	if err != nil {
		return nil, fmt.Errorf("GetChannels: %w", err)
	}
	for _, ch := range resp.Channels {
		if ch.Service == service {
			return &Channel{ID: ch.Id, Name: ch.Name, Service: ch.Service}, nil
		}
	}
	return nil, fmt.Errorf("no channel for service %q in org %s", service, c.orgID)
}
```

> **Note on field names:** `resp.Channels`, `ch.Id`, `ch.Name`, `ch.Service` — these are the genqlient-generated names based on the fields selected in `GetChannels`. If your operation selects different fields or uses aliases, update to match the generated struct field names from `generated.go`.

- [ ] **Step 4: Run the test — expect it to pass**

```bash
go test ./internal/buffer/ -run TestFindChannel -v
```

Expected output:
```
=== RUN   TestFindChannel
=== RUN   TestFindChannel/twitter
=== RUN   TestFindChannel/linkedin
=== RUN   TestFindChannel/tiktok
--- PASS: TestFindChannel (0.00s)
PASS
```

- [ ] **Step 5: Commit**

```bash
git add internal/buffer/client.go internal/buffer/client_test.go
git commit -m "buffer: add BufferClient, FindChannel, and unit tests"
```

---

### Task 7: CreateIdea with mutation response handling

**Files:**
- Modify: `internal/buffer/client.go`, `internal/buffer/client_test.go`

> Before writing, inspect `generated.go` for the `CreateIdeaResponse` type. Look at the type of the field corresponding to `createIdea` in the response. Then follow Pattern A or B below.

- [ ] **Step 1: Determine which response pattern applies**

Open `internal/buffer/generated.go` and search for `CreateIdeaResponse`. Identify:

- **Pattern A (union):** The `CreateIdea` field is typed as an interface (e.g. `CreateIdeaCreateIdea interface { ... }`). There will be two concrete structs implementing it, one for success, one for `MutationError`.
- **Pattern B (plain struct):** The `CreateIdea` field is a plain struct with fields like `Id string` and no interface.

- [ ] **Step 2: Write the failing test**

Add to `internal/buffer/client_test.go`:

**If Pattern A (union):** Replace `CreateIdeaCreateIdeaIdeaActionSuccess` and `CreateIdeaCreateIdeaMutationError` below with the actual concrete type names from `generated.go`.

```go
func TestCreateIdea_Success(t *testing.T) {
	// Fixture matches the success branch of the createIdea union.
	// The JSON field name must match the operation's selection set field name.
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{
			"__typename": "IdeaActionSuccess",
			"idea":       map[string]interface{}{"id": "idea-123"},
		},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	err := c.CreateIdea(context.Background(), "test title", "test text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIdea_Error(t *testing.T) {
	// Fixture matches the MutationError branch.
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{
			"__typename": "MutationError",
			"message":    "invalid organizationId",
			"type":       "NOT_FOUND",
		},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	err := c.CreateIdea(context.Background(), "test title", "test text")
	if err == nil {
		t.Fatal("want error for MutationError branch, got nil")
	}
}
```

**If Pattern B (plain struct):** Use this simpler test instead:

```go
func TestCreateIdea_Success(t *testing.T) {
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{"id": "idea-123"},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	err := c.CreateIdea(context.Background(), "test title", "test text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateIdea_EmptyID(t *testing.T) {
	fixture := map[string]interface{}{
		"createIdea": map[string]interface{}{"id": ""},
	}
	c := &BufferClient{orgID: "test-org", gql: &fakeGQL{data: fixture}}
	err := c.CreateIdea(context.Background(), "test title", "test text")
	if err == nil {
		t.Fatal("want error for empty id, got nil")
	}
}
```

- [ ] **Step 3: Run the test — expect compile error (CreateIdea undefined)**

```bash
go test ./internal/buffer/ -run TestCreateIdea -v
```

Expected: compile error.

- [ ] **Step 4: Add CreateIdea to client.go**

Add to `internal/buffer/client.go`. Replace the concrete type names with the actual names from `generated.go`.

**Pattern A (union):**

```go
// CreateIdea creates a draft idea in the org's Ideas board.
// The mutation returns a union; both the success and MutationError branches are handled.
func (c *BufferClient) CreateIdea(ctx context.Context, title, text string) error {
	resp, err := CreateIdea(ctx, c.gql, c.orgID, title, text)
	if err != nil {
		return fmt.Errorf("CreateIdea: %w", err)
	}
	// Replace the concrete type names below with the actual names from generated.go.
	switch r := resp.CreateIdea.(type) {
	case *CreateIdeaCreateIdeaIdeaActionSuccess:
		_ = r
		return nil
	case *CreateIdeaCreateIdeaMutationError:
		return fmt.Errorf("buffer api error (%s): %s", r.Type, r.Message)
	default:
		return fmt.Errorf("unexpected createIdea response type: %T", resp.CreateIdea)
	}
}
```

**Pattern B (plain struct):**

```go
// CreateIdea creates a draft idea in the org's Ideas board.
func (c *BufferClient) CreateIdea(ctx context.Context, title, text string) error {
	resp, err := CreateIdea(ctx, c.gql, c.orgID, title, text)
	if err != nil {
		return fmt.Errorf("CreateIdea: %w", err)
	}
	if resp.CreateIdea.Id == "" {
		return fmt.Errorf("createIdea returned empty id (check API response)")
	}
	return nil
}
```

- [ ] **Step 5: Run the test — expect pass**

```bash
go test ./internal/buffer/ -v
```

Expected: all tests pass including `TestFindChannel`, `TestCreateIdea_Success`, and `TestCreateIdea_Error` (or `TestCreateIdea_EmptyID`).

- [ ] **Step 6: Commit**

```bash
git add internal/buffer/client.go internal/buffer/client_test.go
git commit -m "buffer: add CreateIdea with mutation response handling"
```

---

### Task 8: Integration smoke test (cmd/probe)

**Files:**
- Create: `cmd/probe/main.go`

- [ ] **Step 1: Write cmd/probe/main.go**

Create `cmd/probe/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davigiroux/bipline/internal/buffer"
)

const (
	orgID   = "6a073bf15d53897094b8a76b"
	service = "twitter"
)

func main() {
	apiKey := os.Getenv("BUFFER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "BUFFER_API_KEY must be set")
		os.Exit(1)
	}

	ctx := context.Background()
	client := buffer.New(apiKey, orgID)

	ch, err := client.FindChannel(ctx, service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FindChannel: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("channel: %s (%s) id=%s\n", ch.Name, ch.Service, ch.ID)

	if err := client.CreateIdea(ctx, "bipline probe", "testing the bipline buffer client — ignore this idea"); err != nil {
		fmt.Fprintf(os.Stderr, "CreateIdea: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("idea created — check Buffer Ideas board")
}
```

- [ ] **Step 2: Build to verify it compiles**

```bash
go build ./cmd/probe/
```

Expected: `./probe` binary produced with no errors.

- [ ] **Step 3: Run against live Buffer API**

```bash
BUFFER_API_KEY=<your-key> ./probe
```

Expected output:
```
channel: devgiroux (twitter) id=6a073cec090476fb99230076
idea created — check Buffer Ideas board
```

Then verify a new idea appears in the Buffer Ideas board at https://publish.buffer.com.

- [ ] **Step 4: Commit**

```bash
git add cmd/probe/main.go
git commit -m "buffer: add probe integration smoke test"
```

---

### Task 9: Package README and PR

**Files:**
- Create: `internal/buffer/README.md`

- [ ] **Step 1: Write the package README**

Create `internal/buffer/README.md`:

```markdown
# internal/buffer

Typed Go client for the Buffer GraphQL API, generated by [genqlient](https://github.com/Khan/genqlient).

## Usage

```go
client := buffer.New(os.Getenv("BUFFER_API_KEY"), orgID)

ch, err := client.FindChannel(ctx, "twitter")
// ch.ID, ch.Name, ch.Service

err = client.CreateIdea(ctx, "title", "text body")
```

## Regenerating bindings

If the Buffer schema changes:

1. `BUFFER_API_KEY=... make schema` — refetch schema.graphql
2. Update `operations.graphql` if needed
3. `make generate` — regenerate generated.go
4. Commit schema.graphql and generated.go in separate commits

## Files

| File | Description |
|------|-------------|
| `genqlient.yaml` | Codegen config (schema, operations, output) |
| `schema.graphql` | Buffer GraphQL schema (introspection JSON) |
| `operations.graphql` | Hand-authored GraphQL operations (Davi owns) |
| `generated.go` | Generated typed client functions — do not edit |
| `client.go` | Thin wrapper: `BufferClient`, `Channel`, auth transport |
```
```

- [ ] **Step 2: Run full test suite**

```bash
go test ./...
```

Expected: all tests pass, no failures.

- [ ] **Step 3: Verify acceptance criteria**

- [ ] `make generate` produces `generated.go` and `go build ./...` is clean
- [ ] `BUFFER_API_KEY=... ./probe` resolves org, finds twitter channel, creates idea
- [ ] `CreateIdea` in `client.go` handles both mutation response branches
- [ ] `internal/buffer/README.md` committed

- [ ] **Step 4: Commit README**

```bash
git add internal/buffer/README.md
git commit -m "buffer: add package README"
```

- [ ] **Step 5: Push branch and open PR**

```bash
git push -u origin main
gh pr create \
  --title "Phase 1: Buffer client" \
  --body "Implements Phase 1 from ROADMAP.md. Adds typed genqlient-backed BufferClient with FindChannel and CreateIdea, unit tests, and integration smoke test (cmd/probe)."
```
