# Story 2.5: Auto-join an open lobby (or create one)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want to automatically get seated in an open lobby for a selected game if one exists,
so that I can start playing with minimal friction.

## Acceptance Criteria

1. **Given** I am authenticated with a valid JWT **When** I request auto-join via `POST /api/v1/lobbies/auto-join` with `game={game_id}` **Then** if an open lobby with an available seat exists, I am joined to that lobby and assigned a seat **And** the response includes the selected `lobby_id` and shareable lobby URL.
2. **Given** there is no open lobby with available seats **When** I call `POST /api/v1/lobbies/auto-join` with `game={game_id}` **Then** a new public lobby is created (2 seats) **And** I am joined to it and assigned a seat **And** the response returns the created `lobby_id` and shareable lobby URL.
3. **Given** I am not authenticated **When** I call `POST /api/v1/lobbies/auto-join` **Then** the request is rejected as unauthorized.

## Tasks / Subtasks

- [ ] Implement `POST /api/v1/lobbies/auto-join` in lobby-service (AC: 1-3)
  - [ ] Require `Authorization: Bearer <token>` and reject missing/invalid tokens
  - [ ] Validate `game` request field is present and matches supported `game_id` values
  - [ ] Search in-memory lobby store for first open public lobby with available seat for this game
  - [ ] If found, join user (reuse join logic from Story 2.4)
  - [ ] If not found, create a new public lobby (reuse create logic from Story 2.3) and seat user
  - [ ] Return lobby details including `lobby_id`, seat counts, and `share_url`
- [ ] Ensure errors and logging align with existing patterns (AC: 1-3)
  - [ ] Use `internal/errorsx` for error formatting
  - [ ] Avoid logging Authorization headers or tokenized URLs
- [ ] Add tests for auto-join flow (AC: 1-3)
  - [ ] Unauthorized request returns 401 + `{error:"forbidden"}`
  - [ ] Missing/invalid `game` returns 400 + `{error:"validation_failed"}`
  - [ ] When open lobby exists, user is joined and counts updated
  - [ ] When no open lobby exists, lobby is created and user seated

## Dev Notes

- Source of truth: `/_bmad-output/planning-artifacts/epics.md` (Story 2.5).
- Prefer existing in-memory store + join/create logic to avoid duplication.
- `max_seats` fixed to 2 in MVP.
- Response JSON must be `snake_case` and include `share_url`.

### Project Structure Notes

- Expected service entrypoint: `cmd/lobby-service/main.go`.
- Suggested handler package: `internal/lobby_service/*`.
- Shared helpers in `internal/auth`, `internal/httpx`, `internal/errorsx`.

### References

- Story 2.5 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-5-Auto-join-an-open-lobby-or-create-one`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX expectations for auto-join: `/_bmad-output/planning-artifacts/ux-design-specification.md#Lobby-Table`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-5-auto-join-an-open-lobby-or-create-one.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.5: Auto-join an open lobby (or create one).
- FRs: FR7 (auto-join) + FR8 (create when none) + supports FR6 (lobby list).

### Business Context

- Auto-join is a core speed-to-seat mechanism; must be fast and reliable.
- Default behavior is “join if possible, otherwise create.”

### UX Context

- Single action: user taps Auto-Join and is seated.
- If a lobby is created, user should see it immediately and be able to invite others.

## Technical Requirements

- REST endpoint: `POST /api/v1/lobbies/auto-join` (authenticated).
- Request must include `game` (string). Reject missing/invalid with 400 `{error:"validation_failed"}`.
- Allowed `game_id` values: `connect4`, `draughts_10x10`.
- Search order: first open **public** lobby with available seat for the game.
- If found, join user using same logic as `/join` (idempotent per user).
- If not found, create a new public lobby with `max_seats=2`, seat creator, and return it.
- Response fields (snake_case):
  - `lobby_id`
  - `game_id`
  - `is_public`
  - `max_seats`
  - `occupied_seats`
  - `seats_available`
  - `created_at`
  - `share_url`
- Data source: in-memory store shared with create/list/join.
- Logging: structured JSON; do not log Authorization headers or tokenized URLs.

## Architecture Compliance

- Base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- JSON naming: `snake_case`.
- Auth: JWT via `Authorization: Bearer <token>`; missing/invalid token is 401 `{error:"forbidden"}`.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.
- JWT verification via `internal/auth` (HS256 signer/verifier).

## File Structure Requirements

- Service entry: `cmd/lobby-service/main.go`.
- Handlers and store: `internal/lobby_service/*`.
- Tests: `test/integration/` (lobby-service HTTP tests).

## Testing Requirements

- Unauthorized request returns 401 + `{error:"forbidden"}`.
- Missing/invalid `game` returns 400 + `{error:"validation_failed"}`.
- When open lobby exists, user is joined and counts updated.
- When no open lobby exists, lobby is created and user seated.
- Ensure search ignores full lobbies.

## Previous Story Intelligence

- Story 2.3 creates public lobbies with `max_seats=2`.
- Story 2.4 join is idempotent and updates seat counts.
- Story 2.2 list uses seat availability fields; keep consistent.

## Git Intelligence Summary

- Recent commits focus on auth-service; lobby-service is still a skeleton.
- Reuse auth middleware patterns from `internal/auth_service/*`.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
