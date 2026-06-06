# DECISION LOG - bipline

Newest first. One entry per non-trivial decision. Keep it short: what, why, what you rejected.

Template:
```
## YYYY-MM-DD - <decision>
Context: <the situation forcing a choice>
Decision: <what you chose>
Alternatives: <what you rejected and why>
```

---

## 2026-06-?? - Seeded from planning <!-- Davi: set the real date -->

### Go, not TypeScript, for the whole tool
Context: the Buffer client, generator, and CLI could be TS or Go.
Decision: Go end to end.
Alternatives: TS (closer to a future React UI) rejected for v0.1. Go is the stronger differentiator to show, and the typed GraphQL client story is clean with genqlient.

### GitHub Action as the trigger, not a daemon or tray
Context: needed an event source for "something shipped."
Decision: a workflow firing on `release.published` runs the CLI.
Alternatives: a long-running daemon or macOS tray (the babygitter shape) rejected. The event we care about is already a CI event, so no infra or token-holding service is needed.

### Buffer is the review surface
Context: build-in-public posts must be reviewed before publishing.
Decision: the tool creates an unscheduled draft, review and scheduling happen inside Buffer.
Alternatives: a custom review UI rejected for v0.1 as unnecessary scope.

### Drafts only, never auto-publish
Context: automate the operational layer, not the judgment.
Decision: bipline never publishes or schedules, it drafts. The human gate stays.

---

## To be decided (gaps)
- [ ] Channel draft vs Idea as the review bucket (resolve in Phase 0)
- [ ] Idempotency mechanism: committed marker file vs Action cache vs querying Buffer for an existing draft (resolve in Phase 4)
- [ ] Notable-commit selection for release drafts (start simple: release notes only)
