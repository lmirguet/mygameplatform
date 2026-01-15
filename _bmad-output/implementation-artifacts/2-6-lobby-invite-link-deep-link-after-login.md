# Story 2.6: Lobby invite link + deep-link after login

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want a shareable lobby invite link that opens the correct lobby after login,
so that I can bring another player into my lobby.

## Acceptance Criteria

1. **Given** I am authenticated and viewing a lobby **When** I request the lobby details via `GET /api/v1/lobbies/{lobby_id}` **Then** the response includes a stable, shareable lobby URL for that lobby.
2. **Given** I am not authenticated **When** I navigate to a lobby URL (deep link) **Then** I am required to log in **And** after successful login, I am redirected back into the same lobby context.
3. **Given** the lobby does not exist (or is invalid) **When** I navigate to its lobby URL **Then** I see a clear “not found” experience with a path back to game/lobby selection.

## Tasks / Subtasks

- [ ] Implement `GET /api/v1/lobbies/{lobby_id}` (AC: 1)
  - [ ] Require `Authorization: Bearer <token>` and reject missing/invalid tokens
  - [ ] Validate `lobby_id` is a UUID; return 400 `{error:"validation_failed"}` if invalid
  - [ ] Look up lobby in in-memory store; return 404 `{error:"not_found"}` if missing
  - [ ] Return lobby details including `lobby_id`, `game_id`, `is_public`, `max_seats`, `occupied_seats`, `seats_available`, `created_at`, and `share_url`
- [ ] Implement invite URL handling in web app (AC: 2-3)
  - [ ] Add route `/lobbies/:lobby_id` that loads lobby context
  - [ ] If not authenticated, redirect to login and preserve intended lobby route
  - [ ] After login, redirect back to intended lobby route
  - [ ] If lobby not found, show “not found” message with CTA back to lobby list
- [ ] Ensure logging and error handling align with existing patterns (AC: 1-3)
  - [ ] Use `internal/errorsx` for API errors
  - [ ] Avoid logging Authorization headers or tokenized URLs
- [ ] Add tests (AC: 1-3)
  - [ ] API: unauthorized returns 401 + `{error:"forbidden"}`
  - [ ] API: invalid UUID returns 400 + `{error:"validation_failed"}`
  - [ ] API: unknown lobby returns 404 + `{error:"not_found"}`
  - [ ] Web: deep link redirects to login when unauthenticated, then back to lobby
  - [ ] Web: invalid lobby route shows not found + CTA

## Dev Notes

- Source of truth: `/_bmad-output/planning-artifacts/epics.md` (Story 2.6).
- Shareable URL is a stable path (`/lobbies/{lobby_id}`) with no hostnames embedded.
- Reuse in-memory lobby store created in Story 2.3.
- Keep response JSON in `snake_case`.

### Project Structure Notes

- API handler: `internal/lobby_service/*`.
- Web routes: `web/src/app/routes.tsx` or equivalent router.
- Use existing auth handling utilities for login/redirect.

### References

- Story 2.6 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-6-Lobby-invite-link--deep-link-after-login`
- API and JSON conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`
- UX deep-link flow: `/_bmad-output/planning-artifacts/ux-design-specification.md#Landing--Invite-→-Join-Lobby-→-Play`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-6-lobby-invite-link-deep-link-after-login.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.6: Lobby invite link + deep-link after login.
- FRs: FR14 (copy invite link), FR15 (deep-link after login), FR19 (shareable URL).

### Business Context

- Invites are a core growth loop; deep-link reliability matters for conversion.
- Users must land in the correct lobby after auth without re-navigation.

### UX Context

- Deep links should be “just works”: minimal friction, clear errors if invalid.
- “Not found” must offer a path back to lobby list.

## Technical Requirements

- REST endpoint: `GET /api/v1/lobbies/{lobby_id}` (authenticated) to retrieve shareable URL and lobby details.
- `lobby_id` must be a UUID; invalid UUID returns 400 `{error:"validation_failed"}`.
- Missing lobby returns 404 `{error:"not_found"}`.
- Response fields (snake_case):
  - `lobby_id`
  - `game_id`
  - `is_public`
  - `max_seats`
  - `occupied_seats`
  - `seats_available`
  - `created_at`
  - `share_url` (stable path like `/lobbies/{lobby_id}`)
- Web app routing:
  - Route `/lobbies/:lobby_id` loads lobby context.
  - If unauthenticated, redirect to login with return path preserved (e.g., query param `next=/lobbies/:lobby_id`).
  - After login, redirect back to the saved path.
  - If lobby not found, show “not found” message + CTA back to lobby list.
- Data source: in-memory lobby store shared with create/list/join.
- Logging: structured JSON; do not log Authorization headers or tokenized URLs.

## Architecture Compliance

- Base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- JSON naming: `snake_case`.
- Auth: JWT via `Authorization: Bearer <token>`; missing/invalid token is 401 `{error:"forbidden"}`.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.
- JWT verification via `internal/auth` (HS256 signer/verifier).
- Frontend routing: `react-router-dom` already selected in architecture.

## File Structure Requirements

- API handlers: `internal/lobby_service/*`.
- Web routes: `web/src/app/routes.tsx` (or equivalent router).
- Auth redirect logic: reuse existing web auth utilities.

## Testing Requirements

- API: unauthorized returns 401 + `{error:"forbidden"}`.
- API: invalid UUID returns 400 + `{error:"validation_failed"}`.
- API: unknown lobby returns 404 + `{error:"not_found"}`.
- Web: unauthenticated deep link redirects to login and back to lobby.
- Web: invalid lobby shows not found + CTA.

## Previous Story Intelligence

- Story 2.3 creates `share_url` as a stable path; reuse the same format.
- Story 2.4 join should be reachable from lobby context after deep-link.
- Story 2.2 list displays lobby availability; align deep-link view with list fields.

## Git Intelligence Summary

- Recent commits focus on auth-service; lobby-service remains a skeleton.
- Reuse auth middleware patterns from `internal/auth_service/*`.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
