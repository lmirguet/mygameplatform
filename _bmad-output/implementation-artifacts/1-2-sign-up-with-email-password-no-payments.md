# Story 1.2: Sign up with email + password (no payments)

Status: review

## Story

As a visitor,
I want to create an account with email + password,
so that I can join lobbies and play games.

## Acceptance Criteria

1. Given I am not authenticated, when I submit a valid email + password to `POST /api/v1/auth/signup`, then a new user account is created and the password is stored as a bcrypt hash (never plaintext) and the response returns an access token (JWT) for the new user.
2. Given I submit an email that is already registered, when I call `POST /api/v1/auth/signup`, then the request is rejected with a clear error response and no duplicate account is created.
3. Given I submit invalid input (e.g., malformed email or too-short password), when I call `POST /api/v1/auth/signup`, then the request is rejected with validation errors.
4. Given the signup endpoint is invoked, when the service logs the request, then it does not log passwords (or password hashes) and it does not include sensitive secrets in logs.

## Tasks / Subtasks

- [x] Preconditions (depends on Story 1.1) (AC: #1)
  - [x] Ensure repo baseline exists (root `go.mod`, `cmd/`, `internal/`, `migrations/`, `web/`) per Story 1.1.
- [x] Define user persistence (AC: #1, #2)
  - [x] Add `users` table via SQL migrations in `migrations/` (only place schema changes happen).
  - [x] Minimum columns (snake_case): `id` (uuid), `email` (unique), `password_hash`, `created_at`, `updated_at`.
  - [x] Add unique index on `email` (e.g., `idx_users__email`).
- [x] Implement `POST /api/v1/auth/signup` in `auth-service` (AC: #1–#4)
  - [x] Validate request body (email format; password min length per policy).
  - [x] Normalize email for uniqueness (define canonicalization strategy; ensure consistent compare/store).
  - [x] On success: create user row; hash password using `golang.org/x/crypto/bcrypt`.
  - [x] Issue JWT access token (minimum claim: `sub` = user id; include TTL from config once defined).
  - [x] Return a success response (JSON) containing at least the access token.
- [x] Handle error cases with consistent formats (AC: #2, #3)
  - [x] Duplicate email → HTTP 409 with `{ "error": "conflict", "message": "Email already registered" }` (or equivalent message).
  - [x] Validation failures → HTTP 400 with `{ "error": "validation_failed", "message": "..." }`.
  - [x] Unexpected errors → HTTP 500 with `{ "error": "internal", "message": "..." }` (avoid leaking internals).
- [x] Logging + security guardrails (AC: #4)
  - [x] Never log passwords/password hashes.
  - [x] Never log JWTs or tokenized URLs (including WS query params, when that arrives).
  - [x] Log errors with structured fields (request_id/trace_id when available).
- [x] Add basic tests (AC: #1–#4)
  - [x] Unit/integration test: successful signup returns token and persists user.
  - [x] Unit/integration test: duplicate signup rejected with 409.
  - [x] Unit/integration test: invalid input rejected with validation_failed.

## Dev Notes

### Architecture Compliance (Non-Negotiable)

- Service boundary: `auth-service` provides `/api/v1` auth endpoints (signup/login) and issues JWTs.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Process Topology; Service Boundaries).
- Password hashing: bcrypt via `golang.org/x/crypto/bcrypt`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Password Hashing).
- Migrations: SQL migrations only, using `golang-migrate/migrate`; schema changes only in `migrations/`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Migrations; Naming Patterns).
- Error format: REST errors always `{error, message}` and error codes are fixed enums (incl. `validation_failed`, `conflict`, `internal`).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (API Response Formats; Error Handling Patterns).
- Naming: DB + JSON are `snake_case`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Naming Patterns; Data Exchange Formats).

### JWT Guardrails

- Architecture flags an early gap: define JWT claims (at minimum `sub` user id), TTL, and signing strategy. Until a dedicated “JWT details” decision story exists, keep signup implementation aligned with architecture defaults and leave TODO hooks in config (do not invent new auth models).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Gap Analysis: JWT details).

### File/Code Placement (Target Layout)

- `cmd/auth-service/main.go` for service entrypoint.
- `internal/auth/` for password + JWT helpers.
- `internal/auth_service/` for HTTP handlers for signup/login and token issuance.
- `internal/errorsx/` for shared error codes.
- `internal/db/` for Postgres helpers and migrations wiring.

Source: `_bmad-output/planning-artifacts/architecture.md` (Complete Project Directory Structure).

### UX Notes (Signup)

- Keep onboarding minimal-field and fast; provide clear validation errors and a forward path.  
  Source: `_bmad-output/planning-artifacts/ux-design-specification.md` (Core UX principles; error patterns).

### References

- `_bmad-output/planning-artifacts/epics.md` → Epic 1 / Story 1.2
- `_bmad-output/planning-artifacts/prd.md` → FR1 (signup, no payment)
- `_bmad-output/planning-artifacts/architecture.md` → auth-service, bcrypt, migrations, error formats, naming conventions

## Dev Agent Record

### Agent Model Used

GPT-5.2 (Codex CLI) — Dev “Amelia”

### Debug Log References

- N/A

### Completion Notes List

- Comprehensive story context created from Epics + Architecture + PRD + UX.
- This story assumes Story 1.1 baseline repo structure exists (or will be completed first).
- Implemented `POST /api/v1/auth/signup` (JSON `{email,password}` → `{access_token}`) with:
  - email validation + normalization (`strings.TrimSpace` + `strings.ToLower`)
  - password policy: minimum length 8
  - bcrypt password hashing
  - HS256 JWT issuance (standard library HMAC-SHA256; optional `exp` when TTL provided)
  - structured request logging without leaking passwords/hashes/JWTs
- Added SQL migrations for `users` table and email unique index.
- Tests:
  - success creates user + returns JWT-like token and ensures no secret leakage to logs
  - duplicate email returns 409 conflict
  - invalid input returns 400 validation_failed

### File List

- `cmd/auth-service/main.go`
- `go.mod`
- `go.sum`
- `internal/auth/jwt_hs256.go`
- `internal/auth/password.go`
- `internal/auth_service/handler.go`
- `internal/auth_service/signup.go`
- `internal/auth_service/signup_test.go`
- `internal/db/postgres.go`
- `internal/db/postgres_users.go`
- `internal/db/users.go`
- `internal/db/uuid.go`
- `internal/errorsx/http.go`
- `internal/httpx/json.go`
- `migrations/000001_create_users_table.down.sql`
- `migrations/000001_create_users_table.up.sql`
- `_bmad-output/implementation-artifacts/1-2-sign-up-with-email-password-no-payments.md`

## Change Log

- 2026-01-13: Implemented `POST /api/v1/auth/signup` + `users` migrations + tests.
