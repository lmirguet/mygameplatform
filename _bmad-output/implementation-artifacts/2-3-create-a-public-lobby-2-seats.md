# Story 2.3: Create a public lobby (2 seats)

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want to create a public lobby for a selected game with exactly 2 player seats,
so that I can invite someone or wait for another player to join.

## Acceptance Criteria

1. **Given** I am authenticated with a valid JWT **When** I create a lobby via `POST /api/v1/lobbies` with `game={game_id}` **Then** a new lobby is created as public/listed **And** the lobby’s max seats is exactly 2 **And** the response returns the created lobby (including its `lobby_id` and shareable URL).
2. **Given** I am not authenticated **When** I call `POST /api/v1/lobbies` **Then** the request is rejected as unauthorized.
3. **Given** I attempt to create a lobby with a seat count other than 2 **When** I call `POST /api/v1/lobbies` **Then** the request is rejected with a validation error.

## Tasks / Subtasks

- [ ] Implement `POST /api/v1/lobbies` handler in lobby-service (AC: 1-3)
  - [ ] Require `Authorization: Bearer <token>` and reject missing/invalid tokens
  - [ ] Validate `game` request field is present and matches supported `game_id` values
  - [ ] Enforce `max_seats` = 2 for MVP; reject any other seat count with validation error
  - [ ] Create a lobby record in an in-memory store for now (no DB in this story)
  - [ ] Generate `lobby_id` (UUID) and set `is_public=true`
  - [ ] Return the created lobby with `lobby_id`, `game_id`, `is_public`, `max_seats`, `occupied_seats`, `seats_available`, `created_at`, and `share_url`
- [ ] Ensure auth + error handling aligns with existing patterns (AC: 1-3)
  - [ ] Use `internal/errorsx` error format and codes
  - [ ] Avoid logging Authorization headers or tokenized URLs
- [ ] Add tests for lobby creation (AC: 1-3)
  - [ ] Unauthorized request returns 401 + `{error:"forbidden"}`
  - [ ] Missing/invalid `game` returns 400 + `{error:"validation_failed"}`
  - [ ] `max_seats` other than 2 returns 400 + `{error:"validation_failed"}`
  - [ ] Happy path returns lobby with `max_seats=2`, `occupied_seats=1` (creator), and valid `share_url`

## Dev Notes

- Source of truth for story requirements: `/_bmad-output/planning-artifacts/epics.md` (Story 2.3).
- Game IDs must match `/api/v1/games` response (`connect4`, `draughts_10x10`).
- Lobby creation is public/listed by default in MVP.
- Shareable URL should be deterministic and align with frontend routing (e.g., `/lobbies/{lobby_id}`); do not hardcode hostnames in server responses.
- Do not add DB schema or persistence in this story; store lobbies in-memory.
- JSON field naming: `snake_case`.

### Project Structure Notes

- Expected service entrypoint: `cmd/lobby-service/main.go`.
- Suggested handler package: `internal/lobby_service/*` (align with `internal/auth_service/*` patterns).
- Shared helpers in `internal/auth`, `internal/httpx`, `internal/errorsx`.

### References

- Story 2.3 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-3-Create-a-public-lobby-2-seats`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX expectations for lobby list/join/create flows: `/_bmad-output/planning-artifacts/ux-design-specification.md#Lobby-Table`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-3-create-a-public-lobby-2-seats.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.3: Create a public lobby (2 seats).
- FRs: FR8 (create lobby) + supports FR7 (join) and FR14/FR19 (shareable URL).

### Business Context

- Creating a lobby is the fast path when no open seat exists and is essential to the 2–4 minute time-to-first-game goal.
- Public/listed lobbies drive liquidity; default to public in MVP.

### UX Context

- Lobby creation is minimal-friction; avoid extra fields.
- Invite link must be easy to copy and aligned with the lobby list UX.

## Technical Requirements

- REST endpoint: `POST /api/v1/lobbies` (authenticated).
- Request must include `game` (string). Reject missing/invalid with 400 `{error:"validation_failed"}`.
- Allowed `game_id` values: `connect4`, `draughts_10x10` (match `/api/v1/games`).
- MVP seat rule: `max_seats` is **always 2**; reject any other requested seat count.
- On create, the creator is immediately seated:
  - `occupied_seats = 1`
  - `seats_available = max_seats - occupied_seats` (must be `1`)
- Response fields (snake_case):
  - `lobby_id` (UUID string)
  - `game_id`
  - `is_public` (true)
  - `max_seats` (2)
  - `occupied_seats` (1)
  - `seats_available` (1)
  - `created_at` (ISO-8601 UTC)
  - `share_url` (stable path like `/lobbies/{lobby_id}`)
- Data source: in-memory store only for this story (no DB schema changes).
- Logging: structured JSON; do not log Authorization headers or tokenized URLs.

## Architecture Compliance

- Base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- JSON and DB naming: `snake_case` everywhere.
- Auth: JWT via `Authorization: Bearer <token>`; missing/invalid token is 401 `{error:"forbidden"}`.
- Do not introduce new error codes or change existing auth conventions.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.
- JWT verification via `internal/auth` (HS256 signer/verifier).

## File Structure Requirements

- Service entry: `cmd/lobby-service/main.go`.
- Handlers: `internal/lobby_service/*` (keep naming aligned with `internal/auth_service/*`).
- Tests: `test/integration/` (lobby-service HTTP tests).

## Testing Requirements

- Unauthorized request returns 401 + `{error:"forbidden"}`.
- Missing/invalid `game` returns 400 + `{error:"validation_failed"}`.
- Seat count other than 2 returns 400 + `{error:"validation_failed"}`.
- Happy path returns lobby with `max_seats=2`, `occupied_seats=1`, `seats_available=1`, and valid `share_url`.
- Share URL is stable and path-only (no hardcoded hostnames).

## Previous Story Intelligence

- Story 2.2 uses an in-memory lobby list with explicit seat availability fields.
- Reuse the same in-memory store for create + list to keep responses consistent.
- Keep `game_id` validation identical to `/api/v1/games` from Story 2.1.

## Git Intelligence Summary

- Recent commits focus on auth-service and Epic 1; lobby-service is still a skeleton.
- No existing lobby handlers or storage to reuse beyond shared `internal/*` helpers.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.
- `golang-migrate/migrate` latest GitHub release is v4.19.0 (released 2025-08-29); keep migrations in `migrations/` when persistence begins in later stories.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
