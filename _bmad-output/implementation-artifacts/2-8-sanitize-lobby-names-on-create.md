# Story 2.8: Sanitize lobby names on create

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want lobby names to be sanitized when I create a lobby,
so that abusive or malformed names are prevented.

## Acceptance Criteria

1. **Given** I am authenticated **When** I create a lobby with a name **Then** the server sanitizes the name (trims whitespace, removes disallowed characters) before storing/returning it.
2. **Given** I provide a name that becomes empty after sanitization **When** I call `POST /api/v1/lobbies` **Then** the request is rejected with a validation error.
3. **Given** I provide an excessively long name **When** I create a lobby **Then** the server truncates or rejects it according to configured limits.

## Tasks / Subtasks

- [ ] Add lobby name support to create endpoint (AC: 1-3)
  - [ ] Extend `POST /api/v1/lobbies` request to accept `name`
  - [ ] Implement sanitization rules (trim, collapse whitespace, allow safe chars)
  - [ ] Enforce min/max length after sanitization
  - [ ] If invalid, return 400 `{error:"validation_failed"}`
  - [ ] Store sanitized name in in-memory lobby record
  - [ ] Return sanitized `name` in create response
- [ ] Update list and detail responses to include sanitized lobby name (AC: 1)
- [ ] Add tests for sanitization (AC: 1-3)
  - [ ] Name with extra whitespace is trimmed/collapsed
  - [ ] Disallowed characters are removed
  - [ ] Empty after sanitization returns validation error
  - [ ] Over max length is rejected or truncated (as specified)

## Dev Notes

- Source of truth: `/_bmad-output/planning-artifacts/epics.md` (Story 2.8).
- Sanitization is part of NFR5: basic protection against lobby abuse.
- Keep name handling consistent across create/list/detail.

### Project Structure Notes

- Lobby handlers: `internal/lobby_service/*`.
- Shared validation helpers can live in `internal/lobby/` or `internal/lobby_service/`.

### References

- Story 2.8 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-8-Sanitize-lobby-names-on-create`
- Architecture error conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-8-sanitize-lobby-names-on-create.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.8: Sanitize lobby names on create.
- FRs: NFR5 (sanitize lobby names; abuse protection).

### Business Context

- Prevent abusive or malformed names from appearing in public lobby lists.
- Keep UX friendly; sanitize where possible instead of rejecting valid user intent.

### UX Context

- Names should look clean and predictable in lobby list table.
- Clear validation error when name is empty or invalid.

## Technical Requirements

- Request: `POST /api/v1/lobbies` accepts optional `name`.
- Sanitization rules (apply in order):
  - Trim leading/trailing whitespace.
  - Collapse internal whitespace to single spaces.
  - Allow characters: letters, numbers, spaces, hyphen, underscore.
  - Remove all other characters.
- Length rules (post-sanitization):
  - Min length: 1 character.
  - Max length: 32 characters.
  - If > 32, truncate to 32 (do not reject).
- If sanitized name becomes empty, return 400 `{error:"validation_failed"}`.
- Store and return sanitized `name` in create/list/detail responses.
- Data source: in-memory lobby store (no DB changes yet).

## Architecture Compliance

- Base path: `/api/v1`.
- Error format: `{ "error": "<code>", "message": "<human>" }` using `internal/errorsx`.
- JSON naming: `snake_case`.
- Do not change existing error codes.

## Library / Framework Requirements

- Go version per `go.mod`: `go 1.25.5`.

## File Structure Requirements

- Lobby handlers: `internal/lobby_service/*`.
- Shared sanitization helper can live in `internal/lobby/` or `internal/lobby_service/`.

## Testing Requirements

- Name is trimmed and whitespace collapsed.
- Disallowed characters are removed.
- Over 32 chars is truncated to 32.
- Empty after sanitization returns 400 `{error:"validation_failed"}`.
- Sanitized name returned in create/list/detail responses.

## Previous Story Intelligence

- Story 2.3 defines create lobby flow and response shape; extend it with `name`.
- Story 2.2 list and Story 2.6 detail should include sanitized `name`.

## Git Intelligence Summary

- Lobby-service code is still minimal; add helper functions alongside new handlers.

## Latest Tech Information

- Go 1.25.5 is the current Go 1.25 patch release (released 2025-12-02); use it for security fixes and compatibility.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
