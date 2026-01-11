---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7]
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/product-brief-mygameplatform-2026-01-11.md
  - _bmad-output/planning-artifacts/product-brief-mygameplatform-2026-01-10_23-13-35.md
  - _bmad-output/planning-artifacts/ux-design-specification.md
  - _bmad-output/planning-artifacts/research/market-multiplayer-browser-game-platform-research-2026-01-11.md
  - _bmad-output/planning-artifacts/research/market-multiplayer-gaming-platform-market-research-2026-01-11.md
workflowType: 'architecture'
project_name: 'mygameplatform'
user_name: 'Laurent'
date: '2026-01-11'
lastStep: 7
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Initialization

Architecture workflow initialized on 2026-01-11 with the following input documents:

- `_bmad-output/planning-artifacts/prd.md`
- `_bmad-output/planning-artifacts/product-brief-mygameplatform-2026-01-11.md`
- `_bmad-output/planning-artifacts/product-brief-mygameplatform-2026-01-10_23-13-35.md`
- `_bmad-output/planning-artifacts/ux-design-specification.md`
- `_bmad-output/planning-artifacts/research/market-multiplayer-browser-game-platform-research-2026-01-11.md`
- `_bmad-output/planning-artifacts/research/market-multiplayer-gaming-platform-market-research-2026-01-11.md`

## Project Context Analysis

### Requirements Overview

**Functional Requirements:**
The PRD groups functional requirements into these architectural areas:

- **Account & Identity**: basic user identity and session handling to support lobby creation/join and ownership.
- **Game Catalog & Selection**: selecting from a small initial catalog (Connect4, Checkers) with room to expand.
- **Lobby Discovery & Joining**: browse/auto-join open lobbies; discoverability and low-friction entry.
- **Hosting & Session Control**: host creates lobbies, starts/ends sessions, shares invites.
- **Gameplay (Authoritative, Rules-Enforced)**: server-authoritative game state and move validation; deterministic rules for shipped games.
- **Invites & Sharing**: shareable URLs and invite flows as a core acquisition mechanic.
- **Rejoin/Recovery (MVP messaging)**: MVP does not require full reconnect, but must handle disconnects gracefully via clear messaging and rapid rejoin flows.
- **Observability & Fair Play Signals**: instrument key product/tech metrics (e.g., time-to-first-game, connection health signals).
- **Web Experience**: browser-based UX with real-time play.
- **Onboarding Speed**: explicit focus on minimizing time-to-first-game.

**Non-Functional Requirements:**
Key NFRs shaping architecture:

- **Performance**: track p95 time-to-first-game and page-to-WebSocket connect time; real-time update expectations.
- **Security/Abuse Controls**: sanitize lobby names; rate limit lobby create/join attempts; basic safety primitives.
- **Scalability**: MVP-scale goals; avoid overbuilding while keeping growth path clear.
- **Accessibility**: UX spec references accessibility needs (WCAG appears), implying baseline a11y considerations in UI/component choices.
- **Integration**: explicitly minimal in MVP (no payments / limited third-party dependencies).

**Scale & Complexity:**

- Primary domain: full-stack web app with real-time multiplayer sessions (browser clients + backend services).
- Complexity level: medium (real-time + session state + abuse controls + observability), despite small initial game count.
- Estimated architectural components: ~9–12 (web client, API/auth, lobby/session service, real-time transport, game/rules engine, persistence layer, telemetry/metrics, moderation primitives, admin tooling/basic ops).

### Technical Constraints & Dependencies

- **Self-hosted deployment** requirement.
- Backend language: **Go**.
- Primary database: **PostgreSQL**.
- Real-time transport indicated as **WebSockets** (PRD).
- MVP avoids external integrations (notably payments).
- Target platforms: desktop + mobile browsers; responsive design required.
- Open question to resolve during architecture decisions: PRD frames **MPA**, while UX references SPA/MPA patterns.

### Cross-Cutting Concerns Identified

- Real-time session lifecycle (create/join/start/end) and state synchronization.
- Deterministic, authoritative game rules enforcement and state integrity.
- Abuse prevention + basic trust/safety primitives (rate limiting, reporting/muting/blocking strategy).
- Observability: product + transport health metrics (TTFG, connection success, disconnect patterns).
- Reliability posture that supports future reconnect/grace windows even if MVP only does messaging.
- Responsive UX and accessibility baseline across devices.

## Starter Template Evaluation

### Primary Technology Domain

Full-stack web application with real-time multiplayer sessions:

- Browser client (mobile + desktop)
- Go backend (HTTP + WebSockets)
- PostgreSQL persistence
- Self-hosted deployment

### Starter Options Considered

1) **Vite + React + TypeScript (build-time Node only)**

- Best fit when you want a Go-only production runtime (frontend can be static assets).
- Vite docs state: Node.js 20.19+ or 22.12+ (templates may require higher). _Source: https://vite.dev/guide/_

2) **Next.js (React)**

- Strong DX and routing, but typically implies a Node runtime in production unless used in a static-only mode.
- Next.js docs list minimum Node.js version 20.9. _Source: https://nextjs.org/docs/app/getting-started/installation_

### Selected Starter: Vite + React + TypeScript

**Rationale for Selection:**

- Matches constraints: Go backend remains the only production server; frontend is static assets.
- Keeps architecture simple for an MVP while supporting real-time gameplay UI.
- Avoids introducing a second server runtime (Node) in prod while still using modern tooling.

**Initialization Command (web client):**

```bash
npm create vite@latest mygameplatform-web -- --template react-ts
```

_Source: https://vite.dev/guide/_

**Architectural Decisions Provided by Starter:**

- Language/runtime (client): TypeScript + modern ESM tooling
- Build tooling: Vite dev server + production build pipeline
- App shape: SPA-style frontend (we can still support “MPA-like” navigation via client-side routes)

**Integration Note (Go backend):**

- Deployment topology: multi-service with a standard ingress terminating TLS on `:443`, routing by hostname to public services. The `auth-service` serves `web/dist` and issues JWTs; the `lobby-service` and `game-service` expose WebSockets behind the ingress.

## Core Architectural Decisions (In Progress)

### Data Architecture

**Database**

- PostgreSQL (primary system of record)

**Migrations**

- Tool: `golang-migrate/migrate` (SQL migrations)
- Policy: “down only when safe; otherwise migrations may be irreversible”
- Version (verified): v4.19.0 (2025-08-29)
  - Source: https://github.com/golang-migrate/migrate/releases

**Game State Persistence (per session): Snapshot-in-DB + Moves**

- Source of truth for “current state”: `game_sessions.current_state` stored as PostgreSQL `JSONB`
- State versioning: include explicit `schema_version` inside `current_state` JSON to support forward-compatible evolution
- Audit/replay trail: append-only `moves` table (one row per move, with timestamp + player + payload)
- Rationale: supports crash recovery + debugging + “watch/replay later” without full event-sourcing complexity

**Core Schema Direction**

- Normalized core tables (MVP): `users`, `lobbies`, `lobby_members`, `game_sessions`, `moves`, `game_results`
- Add: `outbox_events` table to support future async processing (analytics, notifications, moderation signals) without coupling gameplay paths

**Caching (MVP)**

- No Redis initially
- Use in-memory authoritative gameplay + Postgres persistence; revisit Redis when/if we need presence scaling or cross-instance pub/sub

### Authentication & Security

**Authentication UX (MVP)**

- Email + password required (no guest-first flow in MVP)

**Authentication Model (Browser + WebSockets)**

- JWT-based auth
- `auth-service` issues access tokens; other services validate JWTs consistently (same signing configuration/claims)
- WebSockets authenticate via `access_token` query param (see API & Communication Patterns)

**Authorization**

- Simple roles + ownership checks:
  - `user` (default)
  - `admin` (ops/moderation/admin tooling later)
  - Ownership rules: lobby host, seat owner, session participant

**Password Hashing**

- bcrypt (adaptive hash)
  - Source (Go package): https://pkg.go.dev/golang.org/x/crypto/bcrypt

**Rate Limiting (MVP)**

- Per-IP + per-account
  - Apply per-IP to unauthenticated endpoints (signup/login)
  - Apply per-account (and per-IP) to authenticated endpoints (create/join lobby, send invite, websocket connect attempts, etc.)
- Note: with no Redis initially, rate limits are enforced per app instance; revisit if/when you run multiple instances.

**Abuse / Safety Tooling**

- Deferred in MVP (no report/block/mute primitives initially)
- Keep minimal server-side logging hooks for later add-on (so we can backfill basic reporting without redesign).

**TLS Termination**

- TLS termination handled by the ingress/reverse-proxy on `:443` (single edge listener)
- Services are reached via the ingress by hostname; do not run ACME/TLS independently per service

### API & Communication Patterns

**HTTP API Style**

- REST/JSON

**WebSocket Protocol**

- JSON messages (human-readable; fastest iteration for MVP)

**WebSocket Topology**

- Separate WebSockets per feature:
  - `/ws/lobby` (lobby discovery/updates, presence-like signals)
  - `/ws/game` (authoritative gameplay state sync + moves)
- Rationale: isolates traffic patterns and simplifies server-side handlers per domain.

**API Documentation**

- OpenAPI for REST endpoints

**Error Handling Format**

- Simple JSON errors, e.g.:
  - `{ "error": "invalid_credentials", "message": "Invalid email or password" }`
  - `{ "error": "forbidden", "message": "Not lobby host" }`
  - `{ "error": "rate_limited", "message": "Too many requests" }`

**Notes / Implications**

- JWT auth applies to REST and WS.
- WebSocket auth: clients MUST connect with `?access_token=<jwt>` (query param). Servers MUST NOT log full URLs / query strings on WS connect to avoid leaking tokens.
- Define a minimal WS envelope to keep clients stable, e.g.:
  - `{ "type": "lobby.updated", "payload": { ... }, "ts": "..." }`
  - `{ "type": "game.move", "payload": { ... }, "seq": 123 }`

### Frontend Architecture

**Routing**

- React Router (`react-router-dom`)
- Version (verified): 7.12.0
  - Source: https://registry.npmjs.org/react-router-dom/latest

**State Management (MVP)**

- React local state + Context only (no Zustand/Redux in MVP)
- Guidance: keep “server state” in WebSocket streams and derive UI state from WS messages; keep “UI state” local (dialogs, selections, toasts).

**Data Fetching**

- `fetch` + thin wrappers (no TanStack Query in MVP)
- Use REST for non-real-time flows (auth, lobby list bootstrap, game catalog, etc.), and WS for live updates.

**Styling**

- Tailwind CSS
- Version (verified): 4.1.18
  - Source: https://registry.npmjs.org/tailwindcss/latest

**Implications**

- Frontend stores and attaches JWTs for calls to public services (REST + WS); keep token handling consistent across the app.
- Keep routing shallow for MVP: `/` (home), `/lobbies`, `/lobbies/:id`, `/games/:sessionId` (or similar).

### Infrastructure & Deployment

**Packaging**

- Docker images for all components

**Process Topology (Deployable Services)**

- Multiple Go services, deployed as separate containers:
  - `auth-service` (public behind ingress): serves static frontend assets (`web/dist`), provides `/api/v1` auth endpoints (signup/login) and issues JWTs
  - `lobby-service`: lobby discovery/join, lobby membership, invite link resolution, lobby events
  - `game-service`: authoritative game loop, move validation, game session state + persistence (JSONB snapshots + moves)
  - `ingress` (edge): terminates TLS on `:443` and routes by hostname to public services (auth-service, lobby-service, game-service)
  - `postgres`: primary DB
  - `migrate` (job/container): runs `golang-migrate/migrate` at deploy time

**Deployment Target**

- Docker Compose for self-hosting
  - Primary operator command pattern: `docker compose up`
  - Source: https://docs.docker.com/compose/gettingstarted/

**Networking**

- Publicly exposed: `ingress` only (ports 80/443)
- Services are reachable externally only through ingress hostname routing.
- Internal Compose network connects `ingress` ↔ `auth-service`, `lobby-service`, `game-service`, `postgres`.

**TLS Termination**

- TLS in the ingress (single `:443` listener); certs managed by ingress (ACME or provided certs)
- Backend services do not terminate TLS; they accept traffic only from the internal network

**Observability**

- Structured logs only (JSON logs to stdout from all services)
- Correlation: include `request_id` / `trace_id` fields propagated ingress → services

**CI**

- Jenkins builds and tests Go + web client, then builds Docker images
- Jenkins pushes images to your chosen registry (self-hosted or SaaS), then deployment host pulls and runs `docker compose up -d`

**Configuration**

- `.env`-driven configuration for Compose (DB URL, JWT signing config, domain/hostnames, rate limit settings)
- Separate configs for dev vs prod via multiple compose files or compose profiles

## Implementation Patterns & Consistency Rules

### Pattern Categories Defined

**Critical Conflict Points Identified:**
12+ areas where AI agents could make different choices (DB naming, JSON naming, endpoint patterns, WS event shapes, error codes, migration conventions, service boundaries, logging fields, etc.)

### Naming Patterns

**Database Naming Conventions:**

- Tables: plural `snake_case` (e.g., `users`, `lobbies`, `game_sessions`, `outbox_events`)
- Columns: `snake_case` (e.g., `created_at`, `user_id`, `session_id`, `schema_version`)
- PKs: `id` (UUID preferred unless specified otherwise)
- FKs: `<referenced>_id` (e.g., `lobby_id`, `host_user_id`)
- Indexes: `idx_<table>__<cols>` (e.g., `idx_users__email`, `idx_moves__session_id_seq`)
- Timestamps: always `created_at`, `updated_at` (UTC)

**API Naming Conventions (REST):**

- Base path: `/api/v1`
- Resources are plural nouns: `/api/v1/lobbies`, `/api/v1/games`
- IDs are path params: `/api/v1/lobbies/{lobby_id}`
- Query params: `snake_case` (e.g., `page_size`, `include_full`)
- HTTP verbs are standard (GET/POST/PATCH/DELETE)

**WebSocket Naming Conventions:**

- Paths:
  - `/ws/lobby`
  - `/ws/game`
- Event `type`: `snake_case` dotted namespaces:
  - `lobby.created`, `lobby.updated`, `lobby.member_joined`
  - `game.started`, `game.state`, `game.move`, `game.ended`
  - `error.rate_limited`, `error.forbidden`

**Code Naming Conventions:**

- Go packages: use `internal/` for non-public app packages (required)
- Go identifiers: idiomatic Go casing (PascalCase exported, camelCase unexported)
- File names: `snake_case.go` for Go, `kebab-case.tsx` for React components (or keep Vite defaults consistently)
- Services: `auth-service`, `lobby-service`, `game-service`, `ingress` (container/service names)

### Structure Patterns

**Project Organization:**

- Prefer a mono-repo with top-level folders:
  - `cmd/auth-service`, `cmd/lobby-service`, `cmd/game-service`
  - `internal/` shared packages (e.g., `internal/auth`, `internal/db`, `internal/log`, `internal/httpx`, `internal/ws`)
  - `web/` for Vite client (`mygameplatform-web` or standardized to `web/`)
  - `migrations/` for SQL migrations

**File Structure Patterns:**

- All configuration is env-driven; no hardcoded secrets in repo
- Keep Compose files in `deploy/` or root (pick one and stick to it)

### Format Patterns

**API Response Formats:**

- Success: return direct resource JSON (no `{data: ...}` wrapper unless explicitly needed)
- Errors: always:
  - `{ "error": "<error_code>", "message": "<human_readable>" }`
- Error codes are fixed string enums:
  - `invalid_credentials`, `forbidden`, `not_found`, `validation_failed`, `rate_limited`, `conflict`, `internal`

**Data Exchange Formats:**

- JSON fields: `snake_case` everywhere (REST + WS payloads)
- Timestamps: ISO-8601 strings in UTC (e.g., `2026-01-11T14:55:00Z`)
- IDs: UUID strings

### Communication Patterns

**WebSocket Envelope (MANDATORY):**
Every WS message MUST use:

```json
{ "type": "lobby.updated", "payload": { ... }, "ts": "2026-01-11T00:00:00Z", "seq": 123 }
```

- `type`: string
- `payload`: object (can be empty `{}`)
- `ts`: ISO-8601 UTC time
- `seq`: monotonically increasing per-connection stream (or per session stream; choose once and keep consistent)

**State Management Patterns (Client):**

- Treat WS streams as the source of truth for “server state”
- UI-only state stays local to components/context
- JWT handling MUST be consistent (single place in frontend to store/refresh/attach tokens); do not invent per-page token logic

### Process Patterns

**Error Handling Patterns:**

- REST: non-2xx responses always include `{error, message}`
- WS: errors are events with `type` starting `error.` and payload includes at least `{ error, message }`
- Log all server errors with a consistent structured schema (see below)

**Loading State Patterns:**

- Per-route loading indicators for page transitions
- Per-action disabled states for buttons (join/start/move)
- Explicit reconnection UI states (even if MVP does not support session resume)

### Enforcement Guidelines

**All AI Agents MUST:**

- Use `snake_case` for DB and JSON field names everywhere
- Use `/api/v1` for REST endpoints
- Use the WS envelope `{type, payload, ts, seq}` consistently
- Use the error-code enum list; do not invent new codes without updating this section
- Use `golang-migrate/migrate` SQL migrations (no Go migrations)

**Pattern Enforcement:**

- PR checklist: naming + endpoint + WS envelope + error codes
- Add lightweight linters/formatters later; for MVP rely on review + tests

### Pattern Examples

**Good Examples:**

- `POST /api/v1/lobbies` returns `{ "lobby_id": "...", "created_at": "..." }`
- WS message: `{ "type": "game.state", "payload": { "session_id": "...", "schema_version": 1, ... }, "ts": "...", "seq": 42 }`
- Error: `{ "error": "rate_limited", "message": "Too many requests" }`

**Anti-Patterns:**

- Mixing `camelCase` in frontend with `snake_case` in backend payloads
- Returning `{data: ...}` sometimes and raw objects other times
- WS messages without a `type` or with ad-hoc shapes
- Introducing Go migrations or schema diffs outside the migrations folder

## Project Structure & Boundaries

### Complete Project Directory Structure

```
mygameplatform/
├── README.md
├── LICENSE
├── go.mod
├── go.sum
├── Jenkinsfile
├── .gitignore
├── .editorconfig
├── .env.example
├── docs/
│   ├── architecture.md              # points to _bmad-output/planning-artifacts/architecture.md
│   ├── api-openapi.yaml             # REST API spec (source of truth for /api/v1)
│   └── runbooks/
│       ├── local-dev.md
│       └── deploy-compose.md
├── deploy/
│   ├── compose.yaml                 # prod-ish compose
│   ├── compose.override.yaml        # dev overrides (ports, volumes, etc.)
│   ├── secrets/                     # optional (if not using env-only)
│   └── ingress/
│       └── config/                  # ingress config (routes/certs) if not using env-only
├── docker/
│   ├── auth-service/Dockerfile
│   ├── ingress/Dockerfile           # optional (often use off-the-shelf image)
│   ├── lobby-service/Dockerfile
│   ├── game-service/Dockerfile
│   └── web/Dockerfile               # optional if building web in Docker
├── migrations/
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql         # only when safe; otherwise omitted
│   └── ...
├── cmd/
│   ├── auth-service/
│   │   └── main.go
│   ├── lobby-service/
│   │   └── main.go
│   └── game-service/
│       └── main.go
├── internal/
│   ├── auth/                        # password auth (bcrypt), JWT issuance/validation helpers, authz helpers
│   ├── config/                      # env parsing/validation
│   ├── db/                          # pgx/sql boilerplate, tx helpers
│   ├── errorsx/                     # shared error codes + mapping to HTTP/WS
│   ├── httpx/                       # middleware, request_id, JSON helpers, OpenAPI glue
│   ├── log/                         # structured JSON logging helpers
│   ├── rate_limit/                  # per-ip + per-account limiters (instance-local)
│   ├── ws/                          # WS upgrade, origin checks, envelope helpers {type,payload,ts,seq}
│   ├── auth_service/                # auth-service handlers (signup/login, token issuance, serving web/dist)
│   ├── lobby/                       # lobby-service domain (models, repos, handlers)
│   └── game/                        # game-service domain (rules engine, session state, repos)
├── web/
│   ├── package.json
│   ├── package-lock.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   ├── postcss.config.js
│   ├── index.html
│   └── src/
│       ├── main.tsx
│       ├── app/
│       │   ├── routes.tsx           # react-router routes
│       │   └── auth.ts              # token storage + attach to fetch/WS (single source of truth)
│       ├── api/
│       │   ├── client.ts            # fetch wrappers
│       │   └── errors.ts            # maps {error,message} to UI
│       ├── ws/
│       │   ├── lobby_ws.ts          # connects to /ws/lobby, handles envelope
│       │   └── game_ws.ts           # connects to /ws/game, handles envelope
│       ├── pages/
│       │   ├── home/
│       │   ├── lobbies/
│       │   └── game/
│       ├── components/
│       └── styles/
│           └── index.css
└── test/
    ├── integration/
    │   ├── auth_service_test.go
    │   ├── lobby_service_test.go
    │   └── game_service_test.go
    └── fixtures/
        └── ...
```

### Architectural Boundaries

**API Boundaries (Public)**

- Public services are exposed via hostnames routed by ingress:
  - `auth-service` (serves `web/dist` and `/api/v1/...` auth endpoints)
  - `lobby-service` (`/ws/lobby`)
  - `game-service` (`/ws/game`)

**Service Boundaries**

- `auth-service`:
  - serves `web/dist` (static assets)
  - provides `/api/v1` auth endpoints and issues JWTs
- `lobby-service`:
  - lobby creation/join, lobby membership, invites, lobby events
  - writes to Postgres; emits `outbox_events` when relevant
- `game-service`:
  - authoritative move validation + game loop
  - persists `game_sessions.current_state` (JSONB with `schema_version`) and `moves` append-only
  - emits `outbox_events` when relevant

**Data Boundaries**

- Postgres is the shared system of record.
- `migrations/` is the only place schema changes happen.
- JSON and DB use `snake_case` consistently.

### Requirements-to-Structure Mapping (from PRD FR categories)

- **Account & Identity** → `internal/auth`, `auth-service` REST endpoints, `users` tables/migrations
- **Lobby Discovery & Joining / Hosting & Session Control / Invites** → `internal/lobby`, `lobby-service`, related migrations
- **Gameplay (authoritative rules enforced)** → `internal/game`, `game-service`, `game_sessions` + `moves` schema
- **Rejoin/Recovery messaging** → `web/src/ws/*`, token-aware WS reconnect UX patterns
- **Observability & Fair play signals** → `internal/log`, `internal/httpx`, `internal/ws` (request_id/trace_id propagation)
- **Web experience + onboarding speed** → `web/src/pages/*`, `web/src/api/*`, `web/src/ws/*`

### Development / Build / Deploy Structure

**Development**

- `docker compose up` from `deploy/` with dev overrides and mounted volumes for hot reload where applicable.

**Build Process Structure**

- Jenkins:
  - builds Go binaries/images for `auth-service`, `lobby-service`, `game-service`
  - builds `web` and bakes `web/dist` into `auth-service` image (or mounts in dev)
  - runs `migrate` job/container during deployment

**Deployment Structure**

- Compose deploys: `ingress` (80/443) + `auth-service` + `lobby-service` + `game-service` + `postgres`.

## Architecture Validation Results

### Coherence Validation ✅

**Decision Compatibility:**

- Public-edge model is coherent: a single ingress terminates TLS on `:443` and routes by hostname to `auth-service`, `lobby-service`, and `game-service`.
- Auth model is coherent for multi-public services: JWT is used for REST and WebSockets (no cookie-session coupling).
- Web hosting is coherent: `auth-service` serves `web/dist` and provides `/api/v1` auth endpoints.

**Pattern Consistency:**

- Naming conventions are consistent across DB and JSON: `snake_case`.
- REST versioning is consistent: `/api/v1`.
- WS envelope is consistent across lobby/game: `{type, payload, ts, seq}`; WS auth is standardized to `?access_token=<jwt>`.
- Security note captured: services MUST NOT log tokenized WS URLs/query strings.

**Structure Alignment:**

- Mono-repo structure aligns with split deployable services and shared migrations.
- Ingress-as-edge aligns with Docker Compose deployment constraints (single `:443` listener).

### Requirements Coverage Validation ✅

**Functional Requirements Coverage:**

- Account & Identity: covered by `auth-service` issuing JWTs and shared auth helpers.
- Lobby/session flows: covered by `lobby-service` and `/ws/lobby`.
- Gameplay authoritative rules: covered by `game-service` and `/ws/game` with server-side move validation and persisted state.
- Invites/sharing: covered within `lobby-service` and shareable URLs handled by web app + REST.

**Non-Functional Requirements Coverage:**

- Performance: real-time transport specified (WebSockets) and instrumentation paths defined (structured logs with correlation IDs).
- Security: bcrypt for passwords, JWT validation on all services, per-IP + per-account rate limiting, and token-leak logging constraints.
- Scalability: Compose-first MVP with clear path to add centralized rate limiting later (Redis) if needed.

### Implementation Readiness Validation ✅

**Decision Completeness:**

- Core architectural decisions are present (DB/migrations, auth, API/WS conventions, frontend stack, deployment).
- Critical consistency rules are explicit (snake_case, `/api/v1`, WS envelope, error-code enums).

**Structure Completeness:**

- Complete directory structure is defined for mono-repo, services, migrations, and web client.

**Pattern Completeness:**

- Major divergence points are addressed (naming, error codes, WS envelope, migrations).

### Gap Analysis Results

**Critical Gaps (resolve before implementation stories start):**

- **JWT details**: define required claims (at minimum `sub` user id), token TTL, and signing strategy (shared secret vs keypair/JWKS).
- **WS auth leakage**: codify that access tokens in query params must not be logged and should be short-lived; consider one-time WS tokens later if needed.
- **Service-to-service contract**: define whether internal calls are REST or gRPC (and where the OpenAPI specs live if REST). Prevent agents from mixing protocols ad-hoc.
- **Ingress selection/config**: pick the ingress technology (Traefik/Caddy/Nginx) and how hostnames map to services in Compose.

**Important Gaps (should be decided early):**

- **CORS policy** for multi-public services (allowed origins, credentials policy, and whether the web app only talks to same-site subdomains).
- **WS sequencing semantics**: confirm `seq` is per-connection (recommended) and document it explicitly.

### Validation Issues Addressed

- Resolved earlier contradictions by switching from cookie sessions + custom gateway edge to JWT + ingress + `auth-service` hosting `web/dist`.

### Architecture Readiness Assessment

**Overall Status:** READY FOR IMPLEMENTATION (with critical gaps listed above to be resolved as the first implementation decision story)
