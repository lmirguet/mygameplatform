# Story 1.3: Sign in on desktop/mobile (JWT)

Status: done

## Story

As a registered user,
I want to sign in with my email + password,
so that I can access lobbies and gameplay as an authenticated user.

## Acceptance Criteria

1. Given I have a registered account, when I submit valid credentials to `POST /api/v1/auth/login`, then I receive an access token (JWT) and I can use that token to authenticate subsequent API requests.
2. Given I submit invalid credentials, when I call `POST /api/v1/auth/login`, then the request is rejected with a clear error response and the response does not reveal whether the email exists.
3. Given the login endpoint is invoked, when the service logs the request, then it does not log passwords (or password hashes) and it does not include sensitive secrets in logs.

## Tasks / Subtasks

- [x] Preconditions (depends on Story 1.2) (AC: #1)
  - [x] Ensure `users` table exists and passwords are stored as bcrypt hashes.
- [x] Implement `POST /api/v1/auth/login` in `auth-service` (AC: #1–#3)
  - [x] Validate request body (email + password required).
  - [x] Lookup user by email (use same email canonicalization strategy as signup).
  - [x] Verify password using bcrypt compare against stored hash.
  - [x] On success: issue JWT access token (minimum claim: `sub` = user id; TTL from config once defined).
  - [x] Return success response with access token.
- [x] Prevent account enumeration (AC: #2)
  - [x] For invalid email OR invalid password, return the same error code + message (do not reveal existence).
  - [x] Consider constant-time-ish behavior where feasible (at minimum: avoid early returns that create obvious timing gaps).
- [x] Apply rate limiting (AC: #2)
  - [x] Enforce per-IP rate limiting on unauthenticated login attempts (instance-local in MVP).
  - [x] Use standard error format for rate-limited responses.
- [x] Logging + security guardrails (AC: #3)
  - [x] Never log passwords/password hashes.
  - [x] Never log JWTs or tokenized URLs.
  - [x] Log failures in a way that supports troubleshooting without leaking credentials.
- [x] Add tests (AC: #1–#3)
  - [x] Successful login returns token for known user.
  - [x] Invalid credentials returns `invalid_credentials` and does not indicate whether email exists.
  - [x] Rate limiting returns `{ "error": "rate_limited", "message": "Too many requests" }` (or consistent equivalent).

## Dev Notes

### Architecture Compliance (Non-Negotiable)

- `auth-service` exposes `/api/v1` auth endpoints and issues JWTs.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Process Topology; Service Boundaries).
- Password hashing: bcrypt (`golang.org/x/crypto/bcrypt`).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Password Hashing).
- Rate limiting: per-IP applies to unauthenticated endpoints (signup/login); instance-local without Redis in MVP.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Rate Limiting (MVP)).
- Error format: REST non-2xx responses include `{error, message}`; use fixed error codes (`invalid_credentials`, `rate_limited`, `internal`, etc.).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (API Response Formats; Error Handling Format; Error codes enum).

### JWT Guardrails

- JWT details (claims beyond `sub`, TTL, signing strategy) are explicitly called out as an early architecture gap. Don’t invent a bespoke auth model in this story; keep token issuance minimal and configurable.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Gap Analysis: JWT details).

### File/Code Placement (Target Layout)

- `cmd/auth-service/main.go`
- `internal/auth/` (password + JWT helpers)
- `internal/auth_service/` (login handler)
- `internal/errorsx/` (error codes + mapping)

Source: `_bmad-output/planning-artifacts/architecture.md` (Complete Project Directory Structure).

### UX Notes (Login)

- Minimal-field onboarding; clear errors; avoid confusing states.  
  Source: `_bmad-output/planning-artifacts/ux-design-specification.md` (UX principles; error patterns).

### References

- `_bmad-output/planning-artifacts/epics.md` → Epic 1 / Story 1.3
- `_bmad-output/planning-artifacts/prd.md` → FR2 (sign in on desktop/mobile)
- `_bmad-output/planning-artifacts/architecture.md` → auth-service, JWT, bcrypt, rate limiting, error formats

## Dev Agent Record

### Agent Model Used

GPT-5.2 (Codex CLI) — SM “Bob”

### Debug Log References

- N/A

### Completion Notes List

- Comprehensive story context created; includes anti-enumeration requirement and rate-limiting note.
- Implemented login handler with email normalization, bcrypt verification, JWT issuance, and non-enumeration error handling.
- Added per-IP rate limiting for login and ensured sensitive values are not logged.
- Tests added for success, invalid credentials parity, and rate limiting.
- Test run: `go test ./cmd/... ./internal/...` (note: `go test ./...` is blocked by local `pkg/` module cache outside module).
- Code review fixes: trust proxy headers for rate limiting, env-configured rate limit settings, limiter cleanup for expired buckets.
- Tests added for proxy-aware rate limiting; updated `.gitignore` to ignore local `pkg/` module cache.
- `README.md` was already modified in the working tree; not changed during this review, but listed for completeness.

### File List

- `_bmad-output/implementation-artifacts/1-3-sign-in-on-desktopmobile-jwt.md`
- `.gitignore`
- `README.md`
- `cmd/auth-service/main.go`
- `internal/auth_service/handler.go`
- `internal/auth_service/login.go`
- `internal/auth_service/login_test.go`
- `internal/auth_service/rate_limiter.go`
- `internal/auth_service/signup_test.go`
- `internal/db/postgres_users.go`
- `internal/db/users.go`
