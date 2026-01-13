# Story 1.4: View/edit basic profile (username, optional avatar)

Status: ready-for-dev

## Story

As a signed-in user,
I want to view and edit my basic profile (username, optional avatar),
so that my identity is visible and consistent in lobbies and games.

## Acceptance Criteria

1. Given I am authenticated with a valid JWT, when I request my profile via `GET /api/v1/me`, then I receive my current profile fields (at least username and avatar if set).
2. Given I am authenticated with a valid JWT, when I update my profile via `PATCH /api/v1/me` with a new username (and optionally avatar), then my profile is updated and persisted and a subsequent `GET /api/v1/me` reflects the changes.
3. Given I submit invalid profile data (e.g., empty username), when I call `PATCH /api/v1/me`, then the request is rejected with validation errors.

## Tasks / Subtasks

- [ ] Preconditions (depends on Story 1.2 / 1.3) (AC: #1–#3)
  - [ ] Ensure JWT auth middleware exists for REST endpoints.
  - [ ] Ensure `users` table exists and includes a stable user id used for `sub` claim.
- [ ] Extend user model/storage for profile fields (AC: #1–#2)
  - [ ] Add `username` column (required) to `users` table (migration in `migrations/`).
  - [ ] Add `avatar_url` column (nullable) (or equivalent) to `users` table.
  - [ ] Define username validation rules (min/max length, allowed chars, uniqueness policy if required).
- [ ] Implement `GET /api/v1/me` (AC: #1)
  - [ ] Authenticate via JWT and derive user id from token `sub`.
  - [ ] Return profile JSON with at least `{ "username": "...", "avatar_url": "..." }` (snake_case).
- [ ] Implement `PATCH /api/v1/me` (AC: #2–#3)
  - [ ] Authenticate via JWT and derive user id from token `sub`.
  - [ ] Validate payload (reject empty username; enforce validation policy).
  - [ ] Persist updates and return updated profile (or 204 + subsequent GET reflects changes; choose once and keep consistent).
  - [ ] Ensure validation failures return `{ "error": "validation_failed", "message": "..." }`.
- [ ] Add tests (AC: #1–#3)
  - [ ] `GET /api/v1/me` returns 401 when unauthenticated.
  - [ ] `GET /api/v1/me` returns current profile for authenticated user.
  - [ ] `PATCH /api/v1/me` updates username/avatar and subsequent GET reflects changes.
  - [ ] `PATCH /api/v1/me` rejects invalid username with `validation_failed`.

## Dev Notes

### Architecture Compliance (Non-Negotiable)

- REST endpoints are `/api/v1/...` and JSON is `snake_case`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (API Naming Conventions; Data Exchange Formats).
- Error format for REST failures: `{ "error": "<error_code>", "message": "<human_readable>" }`; use fixed enums (incl. `validation_failed`, `forbidden`, `internal`).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (API Response Formats; Error codes enum).
- Account & identity lives in `auth-service` with shared helpers under `internal/auth`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Requirements-to-Structure Mapping; Project Directory Structure).

### JWT Guardrails

- Use `sub` as the user id claim (minimum required per architecture gap list). Avoid inventing additional claims without a dedicated decision.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Gap Analysis: JWT details).

### File/Code Placement (Target Layout)

- `internal/auth_service/` for `/api/v1/me` handlers.
- `internal/auth/` for JWT validation helpers / middleware.
- `internal/db/` for user repository queries/updates.
- `migrations/` for schema changes (only location for schema changes).

Source: `_bmad-output/planning-artifacts/architecture.md` (Complete Project Directory Structure; Migrations policy).

### UX Notes (Profile)

- Keep profile minimal: username + optional avatar; validation errors should be clear and fail-forward.  
  Source: `_bmad-output/planning-artifacts/ux-design-specification.md` (principles: clarity, minimal fields, friendly errors).

### References

- `_bmad-output/planning-artifacts/epics.md` → Epic 1 / Story 1.4
- `_bmad-output/planning-artifacts/prd.md` → FR3
- `_bmad-output/planning-artifacts/architecture.md` → auth model, naming, error format, structure

## Dev Agent Record

### Agent Model Used

GPT-5.2 (Codex CLI) — SM “Bob”

### Debug Log References

- N/A

### Completion Notes List

- Comprehensive story context created for profile read/update endpoints.

### File List

- `_bmad-output/implementation-artifacts/1-4-viewedit-basic-profile-username-optional-avatar.md`
