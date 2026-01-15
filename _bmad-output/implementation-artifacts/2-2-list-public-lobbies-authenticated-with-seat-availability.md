# Story 2.2: List public lobbies (authenticated) with seat availability

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want to view the list of public lobbies for a selected game and see if seats are available,
so that I can quickly join an open game.

## Acceptance Criteria

1. **Given** I am authenticated with a valid JWT **When** I request the lobby list via `GET /api/v1/lobbies?game={game_id}` **Then** I receive a list of public lobbies for that game **And** each lobby includes enough information to determine seat availability (e.g., max seats = 2, current occupied seats).
2. **Given** I am not authenticated **When** I call `GET /api/v1/lobbies?game={game_id}` **Then** the request is rejected as unauthorized.
3. **Given** there are no open lobbies **When** I call `GET /api/v1/lobbies?game={game_id}` **Then** I receive an empty list (not an error).

## Tasks / Subtasks

- [ ] Implement `/api/v1/lobbies` list handler in lobby-service (AC: 1-3)
  - [ ] Validate `game` query param is present and matches supported `game_id` values
  - [ ] Use in-memory stub data source for now (no DB persistence in this story)
  - [ ] Return only public lobbies for the requested game
  - [ ] Include seat availability fields in response with explicit calculation
  - [ ] Return empty list when no lobbies
- [ ] Ensure auth + error handling consistent with auth-service (AC: 1-3)
  - [ ] Require `Authorization: Bearer <token>`; return 401 `{error:"forbidden"}` on missing/invalid token
  - [ ] Use `internal/errorsx` for error formatting
  - [ ] Ensure request logging does not log Authorization header or full URLs containing tokens
- [ ] Add tests for lobby list (AC: 1-3)
  - [ ] Unauthorized request returns 401 + `{error:"forbidden"}`
  - [ ] Missing/invalid `game` param returns 400 + `{error:"validation_failed"}`
  - [ ] No lobbies returns `{ "lobbies": [] }`
  - [ ] Happy path returns lobbies with seat availability fields

## Dev Notes

- Source of truth for story requirements is in `/_bmad-output/planning-artifacts/epics.md` (Story 2.2).
- Use the same `game_id` values returned by `GET /api/v1/games` (Story 2.1) to validate query params.
- Keep API responses consistent with architecture rules: `/api/v1`, `snake_case`, error codes from `internal/errorsx`.
- Boundary: this story is **read-only list** with an **in-memory stub**; do not add DB schema or persistence yet.
- Suggested response shape (example):

```json
{
  "lobbies": [
    {
      "lobby_id": "<uuid>",
      "game_id": "connect4",
      "is_public": true,
      "max_seats": 2,
      "occupied_seats": 1,
      "seats_available": 1,
      "created_at": "2026-01-15T12:00:00Z"
    }
  ]
}
```

### Project Structure Notes

- Expected service entrypoint: `cmd/lobby-service/main.go`.
- Suggested handler package: `internal/lobby_service/*` (aligns with current `internal/auth_service/*` pattern).
- Keep shared helpers in `internal/auth`, `internal/httpx`, `internal/errorsx`.

### References

- Story 2.2 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-2-List-public-lobbies-authenticated-with-seat-availability`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX expectation for lobby list table: `/_bmad-output/planning-artifacts/ux-design-specification.md#Lobby-Table`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-2-list-public-lobbies-authenticated-with-seat-availability.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.2: List public lobbies (authenticated) with seat availability.
- FRs: FR6 (browse open lobbies) + supports FR5 (game selection) and FR7 (join open lobby).

### Business Context

- Fast lobby discovery is core to the 2–4 minute time-to-first-game goal.
- Keeping list minimal and responsive supports lobby liquidity and user confidence.

### UX Context

- Lobby list is displayed as a responsive table; seat availability must be obvious at a glance.
- Empty state should be a friendly “no lobbies” list (not an error).

## Technical Requirements

- REST endpoint: `GET /api/v1/lobbies?game={game_id}` (authenticated).
- `game` query param required; return 400 `{error:"validation_failed"}` if missing/invalid.
- Supported `game_id` values (for now): `connect4`, `draughts_10x10` (must match `/api/v1/games`).
- Only include public lobbies for the requested game.
- Include seat availability fields: `max_seats`, `occupied_seats`, `seats_available`.
- Seat availability formula: `seats_available = max_seats - occupied_seats` (must be >= 0).
- Response contract (types): `lobby_id` UUID string, `game_id` string, `is_public` bool, `max_seats` int, `occupied_seats` int, `seats_available` int, `created_at` ISO-8601 UTC string.
- MVP constraint: lobbies are **2 seats max**; ensure `max_seats` is always `2` for now.
- Performance: response should be fast and small (single-page list; no pagination in this story).
- JSON field naming: `snake_case`.
- Do not log JWTs or full tokenized URLs.
- Rate-limit/cors: follow auth-service conventions (per-IP rate limit where configured, same-site CORS policy if applicable).
- Data source boundary: no DB schema or persistence in this story; schema comes later in Epic 2/3 work.

## Architecture Compliance

- REST base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using codes from `internal/errorsx`.
- Logging: structured JSON logs; avoid sensitive values.
- Service boundary: lobby-service serves this endpoint behind ingress; do not change ingress hostnames or routing in this story.
- Regression guardrail: reuse auth middleware/error handling patterns from `internal/auth_service/*`; do not introduce new error codes.
- Deployment/integration: lobby-service is behind ingress and should remain internal; no new external ports or compose changes in this story.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.
- JWT verification: `internal/auth` (HS256 signer/verifier).

## File Structure Requirements

- New or updated files should live under:
  - `cmd/lobby-service/main.go`
  - `internal/lobby_service/` (handlers, middleware)
  - `test/integration/` (new lobby-service tests)

## Testing Requirements

- Add tests to cover:
  - Unauthorized request returns 401 + `{error:"forbidden"}`.
  - Missing/invalid `game` param returns 400 + `{error:"validation_failed"}`.
  - No lobbies returns `{ "lobbies": [] }`.
  - Happy path includes seat availability fields.
  - `seats_available` computed correctly from `max_seats` and `occupied_seats`.
  - `max_seats` fixed to `2` for MVP lobbies.

## Previous Story Intelligence

- Story 2.1 defines the games list and `game_id` values; reuse those IDs for validation and responses.
- Lobby-service skeleton and auth middleware (if implemented in 2.1) should be reused, not recreated.

## Response Examples

### Empty list

```json
{ "lobbies": [] }
```

### Validation error (missing/invalid `game`)

```json
{ "error": "validation_failed", "message": "game: invalid" }
```

## Definition of Done

- Endpoint responds with correct schema, seat fields, and max seats = 2.
- Unauthorized/missing `game` cases return correct error codes.
- Empty list returns `{ "lobbies": [] }`.
- Tests covering the above are in place.

## Git Intelligence Summary

- Recent commits are focused on Epic 1 and planning artifacts; no lobby-service implementation yet.

## Latest Tech Information

- **Go 1.25 release notes:** Go 1.25 generates DWARF v5 debug info, reducing debug size and link time for large binaries. citeturn1search0
- **React Router DOM:** `react-router-dom` latest version shown on npm is `7.8.2`. citeturn0search1
- **Tailwind CSS:** `tailwindcss` latest version shown on npm is `4.1.12`. citeturn0search0
- **Vite:** Vite docs list Node.js requirement as 20.19+ or 22.12+. citeturn2search2

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
