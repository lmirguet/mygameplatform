# Story 1.1: Set up initial project from starter template

Status: done

## Story

As a developer,
I want to initialize the project using the selected starter template and baseline repo structure,
so that subsequent stories can be implemented efficiently and consistently.

## Acceptance Criteria

1. Given the architecture specifies a Vite + React + TypeScript starter for the web client, when the web project is initialized, then the repository contains a working Vite React TypeScript app (build and dev scripts runnable).
2. And the web app can be served as static assets (`web/dist`) in production.
3. Given the project repository is initialized, when basic build/test checks configured for the repo are run, then they pass.

## Tasks / Subtasks

- [x] Create baseline mono-repo structure aligned to architecture (AC: #1, #2)
  - [x] Create root `go.mod` / `go.sum` (module name `mygameplatform`)
  - [x] Add baseline directories: `cmd/`, `internal/`, `migrations/`, `deploy/`, `docker/`, `docs/`, `test/`
  - [x] Add placeholders as needed so empty dirs are tracked (or add minimal starter files where appropriate)
- [x] Initialize the web client using the approved starter (AC: #1, #2)
  - [x] Create the Vite + React + TypeScript app and place it under `web/` (not a sibling repo)
  - [x] Ensure `npm install` + `npm run dev` works in `web/`
  - [x] Ensure `npm run build` outputs to `web/dist` and can be served as static assets
- [x] Add minimal repo-level “first run” documentation (AC: #3)
  - [x] `README.md` includes: prerequisites, dev commands (web + Go), and where the web build artifacts live (`web/dist`)
  - [x] Add `.env.example` placeholder (values to be filled in later stories)
- [x] Add baseline checks (AC: #3)
  - [x] Confirm `go test ./...` runs successfully (even if only placeholder packages exist initially)
  - [x] Confirm `npm run build` in `web/` succeeds

## Dev Notes

### Architecture Constraints (Do Not Deviate)

- Starter: **Vite + React + TypeScript**. Initialize with `npm create vite@latest mygameplatform-web -- --template react-ts` and then integrate into `web/`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Selected Starter: Vite + React + TypeScript; init command).
- Production web hosting model: `auth-service` serves static frontend assets from `web/dist`.  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Integration Note + service boundaries + repo layout).
- Monorepo directory structure is predefined; use it as the target layout (root `go.mod`, `cmd/*-service`, `internal/*`, `web/`, etc.).  
  Source: `_bmad-output/planning-artifacts/architecture.md` (Complete Project Directory Structure).

### Project Structure Notes

- Current repo is missing `go.mod` and `web/`; this story establishes the baseline structure so later stories can add auth/lobby/game services without churn.
- Keep naming consistent with architecture: service names `auth-service`, `lobby-service`, `game-service`; Go packages under `internal/`; migrations only in `migrations/`.

### Testing Standards (Minimum)

- `go test ./...` at repo root
- `npm run build` under `web/`

### References

- `_bmad-output/planning-artifacts/epics.md` → Epic 1 / Story 1.1 (Acceptance Criteria)
- `_bmad-output/planning-artifacts/architecture.md` → Selected Starter; Project Directory Structure; web/dist serving model
- `_bmad-output/planning-artifacts/prd.md` → MVP is browser-first (desktop/mobile), real-time play later; no payments in MVP

## Dev Agent Record

### Agent Model Used

GPT-5.2 (Codex CLI) — Dev “Amelia”

### Debug Log References

- N/A

### Completion Notes List

- Initialized monorepo baseline (Go module `mygameplatform`, service entrypoints under `cmd/*-service`, tracked empty dirs via `.gitkeep` placeholders).
- Bootstrapped `web/` as Vite + React + TypeScript and verified `npm run build` produces `web/dist`.
- Added repo-level first-run docs + `.env.example`.
- Verified checks:
  - `GOCACHE=$PWD/.go-cache GOMODCACHE=$PWD/.go-modcache go test ./...`
  - `cd web && npm_config_cache=$PWD/../.npm-cache npm run build`
- Environment workaround: local caches (`.npm-cache/`, `.go-cache/`, `.go-modcache/`) to avoid global cache permission issues.

### File List

- `.env.example`
- `.gitignore`
- `README.md`
- `cmd/auth-service/main.go`
- `cmd/game-service/main.go`
- `cmd/lobby-service/main.go`
- `deploy/.gitkeep`
- `deploy/ingress/config/.gitkeep`
- `docker/.gitkeep`
- `docker/auth-service/.gitkeep`
- `docker/game-service/.gitkeep`
- `docker/ingress/.gitkeep`
- `docker/lobby-service/.gitkeep`
- `docker/web/.gitkeep`
- `docs/.gitkeep`
- `docs/runbooks/.gitkeep`
- `go.mod`
- `go.sum`
- `migrations/.gitkeep`
- `test/.gitkeep`
- `test/fixtures/.gitkeep`
- `test/integration/.gitkeep`
- `web/.gitignore`
- `web/README.md`
- `web/eslint.config.js`
- `web/index.html`
- `web/package-lock.json`
- `web/package.json`
- `web/public/vite.svg`
- `web/src/App.css`
- `web/src/App.tsx`
- `web/src/assets/react.svg`
- `web/src/index.css`
- `web/src/main.tsx`
- `web/tsconfig.app.json`
- `web/tsconfig.json`
- `web/tsconfig.node.json`
- `web/vite.config.ts`
- `_bmad-output/implementation-artifacts/1-1-set-up-initial-project-from-starter-template.md`

## Senior Developer Review (AI)

Reviewer: Amelia (DEV agent) on 2026-01-14

Findings (fixed):
- [AC#2] Serve `web/dist` as static assets: `cmd/auth-service/main.go`
- [AC#3] Document baseline checks + static hosting env: `README.md`, `.env.example`
- Add tests for SPA static hosting behavior: `internal/httpx/spa_test.go`

Notes:
- This review workflow’s git-diff checks only detect *uncommitted* changes; repository was clean at review time.

### Change Log

- 2026-01-14: Code review fixes applied (static hosting + docs + tests).
