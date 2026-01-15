# Story 2.4: Join a lobby and claim a seat

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want to join a specific lobby and claim an available seat,
so that I can participate in the upcoming session.

## Acceptance Criteria

1. **Given** I am authenticated with a valid JWT **When** I join a lobby via `POST /api/v1/lobbies/{lobby_id}/join` **Then** I become a lobby member assigned to an available player seat.
2. **Given** the lobby is full (both seats taken) **When** I call `POST /api/v1/lobbies/{lobby_id}/join` **Then** the request is rejected with a clear “lobby full” error.
3. **Given** I already joined the lobby previously **When** I call `POST /api/v1/lobbies/{lobby_id}/join` again **Then** the operation is idempotent (no duplicate membership) **And** I remain assigned to my existing seat.
4. **Given** I am not authenticated **When** I call `POST /api/v1/lobbies/{lobby_id}/join` **Then** the request is rejected as unauthorized.

## Tasks / Subtasks

- [ ] Implement `POST /api/v1/lobbies/{lobby_id}/join` in lobby-service (AC: 1-4)
  - [ ] Require `Authorization: Bearer <token>` and reject missing/invalid tokens
  - [ ] Validate `lobby_id` is present and is a UUID
  - [ ] Look up lobby in in-memory store created in Story 2.3
  - [ ] If lobby not found, return 404 `{error:"not_found"}`
  - [ ] If user already in lobby, return existing seat assignment (idempotent)
  - [ ] If lobby full (occupied_seats == max_seats), return 409 `{error:"conflict"}` with clear message
  - [ ] Otherwise assign next available seat and increment `occupied_seats`
  - [ ] Return lobby membership response with `lobby_id`, `user_id`, `seat_index`, and lobby seat counts
- [ ] Ensure error handling & logging align with existing patterns (AC: 1-4)
  - [ ] Use `internal/errorsx` error format and codes
  - [ ] Avoid logging Authorization headers or tokenized URLs
- [ ] Add tests for join flow (AC: 1-4)
  - [ ] Unauthorized request returns 401 + `{error:"forbidden"}`
  - [ ] Unknown lobby returns 404 + `{error:"not_found"}`
  - [ ] Full lobby returns 409 + `{error:"conflict"}` with “lobby full” message
  - [ ] Idempotent re-join returns same seat assignment and does not increment counts
  - [ ] Happy path assigns seat and increments `occupied_seats`

## Dev Notes

- Source of truth for story requirements: `/_bmad-output/planning-artifacts/epics.md` (Story 2.4).
- Seat count is fixed to 2 in MVP.
- Join must be idempotent per user per lobby.
- Keep response JSON in `snake_case`.
- Use the same in-memory store as Story 2.3 to ensure join/list/create stay consistent.

### Project Structure Notes

- Expected service entrypoint: `cmd/lobby-service/main.go`.
- Suggested handler package: `internal/lobby_service/*`.
- Shared helpers in `internal/auth`, `internal/httpx`, `internal/errorsx`.

### References

- Story 2.4 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-4-Join-a-lobby-and-claim-a-seat`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX expectations for join flow: `/_bmad-output/planning-artifacts/ux-design-specification.md#Lobby-Table`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-4-join-a-lobby-and-claim-a-seat.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.4: Join a lobby and claim a seat.
- FRs: FR7 (join lobby) + supports FR6 (list lobbies) and FR19 (shareable URL entry).

### Business Context

- Fast seat acquisition is critical to the 2–4 minute time-to-first-game goal.
- Join must be reliable and idempotent to prevent duplicate seats.

### UX Context

- Clear error feedback for full lobbies; user should be able to return to list and try another.
- Join flow must be a single action without extra fields.

## Technical Requirements

- REST endpoint: `POST /api/v1/lobbies/{lobby_id}/join` (authenticated).
- `lobby_id` must be a UUID path param; invalid UUID returns 400 `{error:"validation_failed"}`.
- If lobby not found, return 404 `{error:"not_found"}`.
- If lobby is full, return 409 `{error:"conflict"}` with message “lobby full”.
- Join is idempotent per user per lobby:
  - Repeated join returns same seat assignment.
  - Do not increment seat counts on re-join.
- Seat assignment:
  - `max_seats` fixed to 2 in MVP.
  - Assign lowest available seat index (0 or 1).
  - After successful join: `occupied_seats` increments by 1 and `seats_available` decrements.
- Response fields (snake_case):
  - `lobby_id`
  - `user_id`
  - `seat_index`
  - `max_seats`
  - `occupied_seats`
  - `seats_available`
  - `joined_at` (ISO-8601 UTC)
- Data source: same in-memory lobby store used in Story 2.3.
- Logging: structured JSON; do not log Authorization headers or tokenized URLs.

## Architecture Compliance

- Base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- JSON and DB naming: `snake_case`.
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
- Invalid UUID returns 400 + `{error:"validation_failed"}`.
- Unknown lobby returns 404 + `{error:"not_found"}`.
- Full lobby returns 409 + `{error:"conflict"}` with “lobby full”.
- Idempotent re-join returns same seat index and does not increment counts.
- Happy path assigns seat and updates counts correctly.

## Previous Story Intelligence

- Story 2.3 creates public lobbies with `max_seats=2` and initial `occupied_seats=1`.
- Story 2.2 lists lobbies with explicit seat availability fields.
- Reuse the same in-memory store across create/list/join for consistency.

## Git Intelligence Summary

- Recent commits focus on auth-service; lobby-service remains a skeleton.
- Use existing auth middleware patterns from `internal/auth_service/*`.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
