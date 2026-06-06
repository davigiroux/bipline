# ROADMAP - bipline

Current phase: **Phase 2 - Draft generator** <!-- Davi: keep this in sync with CLAUDE.md -->

Advance only when every box in a phase is checked.

## Phase 0 - Spike (manual, no code) ✅
- [x] Buffer API key works against `https://api.buffer.com`
- [x] Listed organizations, recorded org ID in CLAUDE.md
- [x] Listed channels, recorded service -> channel ID mapping
- [x] Created one draft by hand (curl), confirmed required args
- [x] Decided channel draft vs Idea (record in DECISION-LOG.md)

Notes: <!-- Davi -->

## Phase 1 - Buffer client (Go) ✅
- [x] genqlient generates bindings from the live schema
- [x] Resolves org, finds channel by service, creates a draft
- [x] Both mutation result branches handled (CreateIdeaPayload union, 6 variants)
- [x] Short package README committed

Notes: <!-- Davi -->

## Phase 2 - Draft generator (Go)
- [ ] `generator.Draft(event, voice)` returns a draft string
- [ ] Golden tests on 2-3 real past events pass voice invariants
- [ ] No Buffer network calls in this package

Notes:

## Phase 3 - Event source (Go)
- [ ] Parses a `release.published` fixture into a normalized Event
- [ ] Parses a merged-PR fixture into a normalized Event
- [ ] Unit tested against fixture JSON

Notes:

## Phase 4 - Glue: CLI + Action
- [ ] CLI reads `GITHUB_EVENT_PATH`, drafts, creates a Buffer draft
- [ ] Workflow committed in SafeNudge
- [ ] A test release produces exactly one draft
- [ ] A re-run produces zero new drafts (idempotency works)

Notes:

## Phase 5 - Dogfood
- [ ] bipline repo has the workflow
- [ ] Cutting `v0.1` drafts the announcement post
- [ ] Reviewed in Buffer and published (the first post)

Notes:

## Phase 6 - Optional
- [ ] Merged-PR trigger with label gate
- [ ] Multi-repo hub on a cron
- [ ] Per-platform variants (X / LinkedIn)
- [ ] React review pane (only if Buffer's review UX is not enough)

Notes:
