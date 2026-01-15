# Story 2.1: List games (authenticated)

Status: done

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want to view the available games on the platform,
so that I can choose what to play next.

## Acceptance Criteria

1. **Given** I am authenticated with a valid JWT **When** I request the games list via `GET /api/v1/games` **Then** I receive a list that includes at least Connect4 and Draughts (10x10).
2. **Given** I am not authenticated **When** I call `GET /api/v1/games` **Then** the request is rejected as unauthorized.
3. **Given** the games list endpoint is invoked **When** the service logs the request **Then** it does not log sensitive auth material (e.g., raw JWTs).

## Tasks / Subtasks

- [x] Build lobby-service HTTP server skeleton and routing (AC: 1-3)
  - [x] Create `cmd/lobby-service/main.go` HTTP server similar to auth-service (env config, logger, mux, timeouts)
  - [x] Add `internal/lobby_service/handler.go` (or similar) with router for `/api/v1/games`
  - [x] Implement auth middleware using `internal/auth` verifier and `Authorization: Bearer` header (reuse error format from `internal/errorsx`)
- [x] Implement `/api/v1/games` handler (AC: 1-3)
  - [x] Return JSON with a stable list of games (no DB yet) and basic rules/players info
  - [x] Ensure JSON uses `snake_case` fields and error codes match `errorsx` conventions
  - [x] Ensure request logging does not record Authorization header or tokenized URLs
- [x] Add tests for lobby-service games list (AC: 1-3)
  - [x] Unit test auth middleware and unauthorized response
  - [x] Integration-style test for `/api/v1/games` happy path (valid JWT)

## Dev Notes

- Source of truth for story requirements is in `/_bmad-output/planning-artifacts/epics.md` (Story 2.1).
- Use the existing auth patterns from `internal/auth_service/*` as reference (auth middleware, error formats, logging).
- Keep API responses consistent with architecture rules: `/api/v1`, `snake_case`, error codes from `internal/errorsx`.
- Suggested response shape (example):

```json
{
  "games": [
    {
      "game_id": "connect4",
      "name": "Connect4",
      "min_players": 2,
      "max_players": 2,
      "rules_summary": "7x6 board; gravity drop; 4-in-a-row wins"
    },
    {
      "game_id": "draughts_10x10",
      "name": "Draughts (10x10)",
      "min_players": 2,
      "max_players": 2,
      "rules_summary": "International Draughts: forced captures, max-capture priority, multi-jumps, flying kings"
    }
  ]
}
```

### Project Structure Notes

- Expected service entrypoint: `cmd/lobby-service/main.go`.
- Suggested handler package: `internal/lobby_service/*` (aligns with current `internal/auth_service/*` pattern).
- Keep shared helpers in `internal/auth`, `internal/httpx`, `internal/errorsx`.

### References

- Story 2.1 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-1-List-games-authenticated`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX expectation for game list visibility: `/_bmad-output/planning-artifacts/ux-design-specification.md#Lobby-Table`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Implementation Plan

- Implement lobby-service HTTP server skeleton, routing, and auth middleware first; keep games handler stubbed until task 2 to preserve task order.

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created
- Implemented lobby-service server skeleton with `/api/v1/games` routing + auth middleware; added initial auth-required test. Tests: `go test ./internal/lobby_service`, `go test ./cmd/... ./internal/...` (note: `go test ./...` fails due to `pkg/mod` directory in repo).
- Implemented `/api/v1/games` response with Connect4 + Draughts (10x10) and added catalog + log-safety tests. Tests: `go test ./internal/lobby_service`, `go test ./cmd/... ./internal/...`.
- Completed lobby-service games list test coverage (auth required + catalog happy path). Tests: `go test ./cmd/... ./internal/...`.
- Story 2.1 complete; status moved to review. Tests: `go test ./cmd/... ./internal/...` (note: `go test ./...` still fails due to repo `pkg/mod`).
- Code review fixes: added integration test under `test/integration/` and switched lobby-service logs to JSON. Tests: `go test ./cmd/... ./internal/... ./test/...`.

### File List

- `/_bmad-output/implementation-artifacts/2-1-list-games-authenticated.md`
- `/_bmad-output/implementation-artifacts/sprint-status.yaml`
- `cmd/lobby-service/main.go`
- `internal/lobby_service/auth_middleware.go`
- `internal/lobby_service/games.go`
- `internal/lobby_service/games_test.go`
- `internal/lobby_service/handler.go`
- `internal/lobby_service/handler_test.go`
- `internal/lobby_service/log.go`
- `test/integration/lobby_service_test.go`

### Change Log

- 2026-01-15: Implemented lobby-service games catalog endpoint + auth middleware; added tests and logging safeguards.
- 2026-01-15: Code review fixes: JSON logging + integration test placement.

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.1: List games (authenticated).
- FRs: FR4 (view available games); supports FR5 (game selection into join/host flow).

### Business Context

- MVP games are limited to Connect4 and Draughts (10x10) to maximize lobby liquidity and reduce scope.
- Game list is a core step in the “speed to seat” UX; keep response fast and stable.

### UX Context

- Lobby list and game selection are designed as a fast, readable table with clear CTA (join/create).
- Mobile-first responsive layouts; games list should be simple to render on both desktop and mobile.

## Technical Requirements

- REST endpoint: `GET /api/v1/games` (authenticated).
- Auth: JWT in `Authorization: Bearer <token>`; reject missing/invalid tokens with 401 + `{error:"forbidden"}`.
- Return a stable list including **Connect4** and **Draughts (10x10)** with basic rules/players info.
- Do not log JWTs or full tokenized URLs (especially for WebSocket tokens elsewhere).
- JSON field naming: `snake_case`.

## Architecture Compliance

- REST base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using codes from `internal/errorsx`.
- Logging: structured JSON logs; avoid sensitive values.
- Consistency with existing auth-service patterns for auth + error handling.

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
  - Authorized request returns list including `connect4` and `draughts_10x10`.
  - Response JSON uses `snake_case` fields.

## Previous Story Intelligence

- Previous epic stories (Epic 1) implemented auth + profile in `auth-service`. Reuse auth patterns; no lobby-service code exists yet.

## Git Intelligence Summary

- No lobby-service commits exist; only auth-service and shared helpers are present.

## Latest Tech Information

- **Go runtime:** Go `1.25.5` is the latest published release in the 1.25 line (security fixes included). citeturn0search4turn0search1
- **React Router DOM:** `react-router-dom` latest tag is `7.12.0` (Node >= 20). citeturn2view0
- **Tailwind CSS:** `tailwindcss` latest tag is `4.1.18`. citeturn2view1
- **Vite:** Vite 7 is the current major release; Node.js support is 20.19+ or 22.12+. Latest published Vite package version appears as `7.3.1`. citeturn0search0turn6search5

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **done**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
