# Story 2.9: Rate limit lobby create/join attempts

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want lobby create/join attempts to be rate limited,
so that abuse is controlled and the system remains stable.

## Acceptance Criteria

1. **Given** I am authenticated **When** I attempt to create lobbies rapidly **Then** requests beyond the rate limit are rejected with a clear “rate limited” error.
2. **Given** I am authenticated **When** I attempt to join lobbies rapidly **Then** requests beyond the rate limit are rejected with a clear “rate limited” error.
3. **Given** I am within limits **When** I create/join lobbies **Then** requests succeed as normal.

## Tasks / Subtasks

- [ ] Add rate limiting to lobby create/join endpoints (AC: 1-3)
  - [ ] Apply per-account limiter for `POST /api/v1/lobbies`
  - [ ] Apply per-account limiter for `POST /api/v1/lobbies/{lobby_id}/join`
  - [ ] Apply per-account limiter for `POST /api/v1/lobbies/auto-join`
  - [ ] Return 429 `{error:"rate_limited"}` when exceeded
- [ ] Ensure limiter configuration matches auth-service patterns
  - [ ] Reuse existing rate limiter helper (if available)
  - [ ] Do not log sensitive tokens or full URLs
- [ ] Add tests for rate limiting (AC: 1-3)
  - [ ] Excess create attempts return 429 + `{error:"rate_limited"}`
  - [ ] Excess join attempts return 429 + `{error:"rate_limited"}`
  - [ ] Within limits succeeds

## Dev Notes

- Source of truth: `/_bmad-output/planning-artifacts/epics.md` (Story 2.9).
- Rate limiting is part of NFR5 (abuse protection).
- Use per-account rate limits; per-IP limits are handled elsewhere.

### Project Structure Notes

- Lobby handlers: `internal/lobby_service/*`.
- Rate limit helper likely in `internal/rate_limit` or reuse `auth_service` limiter.

### References

- Story 2.9 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-9-Rate-limit-lobby-createjoin-attempts`
- Architecture security notes: `/_bmad-output/planning-artifacts/architecture.md#Authentication--Security`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-9-rate-limit-lobby-createjoin-attempts.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.9: Rate limit lobby create/join attempts.
- FRs: NFR5 (rate limit create/join; abuse protection).

### Business Context

- Prevent abuse and protect service stability while maintaining low friction.

### UX Context

- When rate limited, show clear, friendly error message and allow retry after cooldown.

## Technical Requirements

- Apply per-account rate limits to:
  - `POST /api/v1/lobbies`
  - `POST /api/v1/lobbies/{lobby_id}/join`
  - `POST /api/v1/lobbies/auto-join`
- On limit exceeded, return 429 `{error:"rate_limited"}` with clear message.
- Rate limit scope: per-account (user id) using JWT `sub`.
- If no user id available, treat as unauthorized (401).
- Default limits (until config exists):
  - Create: 5 requests / minute
  - Join: 10 requests / minute
  - Auto-join: 5 requests / minute
- Logging: structured JSON; do not log Authorization headers or tokenized URLs.

## Architecture Compliance

- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- Use existing rate limiter patterns from auth-service when possible.
- Keep endpoints and JSON naming unchanged.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.

## File Structure Requirements

- Lobby handlers: `internal/lobby_service/*`.
- Rate limit helpers: reuse `internal/rate_limit` or copy pattern from `internal/auth_service/rate_limiter.go`.

## Testing Requirements

- Excess create attempts return 429 `{error:"rate_limited"}`.
- Excess join attempts return 429 `{error:"rate_limited"}`.
- Excess auto-join attempts return 429 `{error:"rate_limited"}`.
- Requests within limits succeed.

## Previous Story Intelligence

- Auth-service already uses a rate limiter; follow the same approach.
- Lobby endpoints from stories 2.3–2.5 should reuse a shared limiter.

## Git Intelligence Summary

- Lobby-service not yet implemented; rate limiters can mirror auth-service patterns.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
