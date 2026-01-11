---
stepsCompleted:
  - step-01-validate-prerequisites
  - step-02-design-epics
  - step-03-create-stories
  - step-04-final-validation
inputDocuments:
  - _bmad-output/planning-artifacts/prd.md
  - _bmad-output/planning-artifacts/architecture.md
  - _bmad-output/planning-artifacts/ux-design-specification.md
---

# mygameplatform - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for mygameplatform, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: A visitor can create an account in one flow without providing payment information.
FR2: A registered user can sign in on desktop or mobile browser.
FR3: A user can view and edit basic profile fields (username, avatar optional) without payment details.
FR4: A user can view available games (Connect4, Checkers) with basic rules/players info.
FR5: A user can pick a game and proceed directly to join/host flow.
FR6: A user can see a list of open lobbies for a selected game with seat availability.
FR7: A user can auto-join an open lobby and take an available seat.
FR8: If no suitable lobby exists, a user can create a new lobby for that game.
FR9: A host can start a session once minimum seats are filled and end the session at any time.
FR10: A host can remove a player from the session before or during play.
FR11: Players in a session can make moves according to game rules; illegal moves are rejected.
FR12: The system updates and shares canonical game state to all session players after each valid move.
FR13: A session concludes with a clear game result (win/lose/draw) per rules for Connect4 and Checkers.
FR21: The system implements Connect4 rules (7x6 board, gravity drop into a column, alternating turns) and rejects illegal moves (invalid column, full column, wrong turn).
FR22: The system detects and communicates Connect4 outcomes (4-in-a-row horizontal/vertical/diagonal, draw when board is full).
FR23: The system implements draughts/checkers gameplay on a 10x10 board (per your chosen ruleset) and rejects illegal moves.
FR24: The system enforces forced captures for draughts/checkers (a capture must be made when available), including backward captures for men.
FR25: The system enforces multi-jump capture sequences for draughts/checkers within a single turn when available (continuation captures required until none remain).
FR26: The system performs kinging for draughts/checkers when a piece reaches the last rank, and kings have “flying” movement/captures (multi-square diagonals) per your ruleset.
FR27: The system detects and communicates draughts/checkers outcomes (win/loss), and supports draws when both players agree to a draw.
FR28: The web client provides clear, game-specific feedback for illegal moves and shows the current turn, last move, and result state for both games.
FR14: A host or player can copy an invite link to a lobby/session for others to join.
FR15: An invited user can follow a link and land directly in the correct game/lobby context after signup/login.
FR16: If a player disconnects, the system communicates loss of session and offers fast path to rejoin a new lobby (no state restore in MVP).
FR17: The system records session events (join, start, moves, end, kicks) for basic troubleshooting and fairness checks.
FR18: The product is usable on recent desktop and mobile browsers with responsive layouts and touch/mouse input.
FR19: Each lobby/session has a shareable URL; users can load it directly to join or host if eligible.
FR20: The flow from landing to sitting in a lobby enables completion within the 2–4 minute target (as tracked in metrics).

### NonFunctional Requirements

NFR1: (Performance) Time-to-first-game p95 ≤ 4 minutes (tracked); page-to-WebSocket connect time monitored.
NFR2: (Performance) Game state updates visible to players within a small real-time window (WebSocket-based); no strict SLO set in MVP.
NFR3: (Security) Accounts require authenticated access; no payments collected in MVP.
NFR4: (Security) Session access via invite link must validate eligibility (seat available, game context).
NFR5: (Security) Basic protection against lobby abuse: sanitize lobby names; rate limit lobby create/join attempts.
NFR6: (Scalability) MVP targets WAU 100 / MAU 200; design to add more lobbies/games without architectural change.
NFR7: (Accessibility) No formal requirement for MVP; responsive layouts and touch/mouse input supported.
NFR8: (Integration) None required in MVP (no external payments or third-party systems).

### Additional Requirements

- Starter template (web client): Vite + React + TypeScript.
- Platform/runtime: Go backend + PostgreSQL; self-hosted deployment.
- Deployment topology: TLS ingress on `:443`; `auth-service` serves `web/dist` and provides `/api/v1/...` auth endpoints; `lobby-service` + `game-service` provide WebSockets behind ingress.
- Authentication model: JWT for REST and WebSockets; email + password (no guest-first flow in MVP); bcrypt password hashing.
- WebSocket auth constraint: clients connect with `?access_token=<jwt>`; services must not log full URLs/query strings to avoid token leakage.
- API conventions: REST base path `/api/v1`; consistent resource naming; consistent error conventions.
- WebSocket envelope: consistent message envelope (e.g., `type` + `payload` + `ts` + `seq`) for real-time events/state.
- Data persistence: `game_sessions.current_state` stored as Postgres `JSONB` and must include a `schema_version`; append-only `moves` table for replay/audit.
- Migrations: SQL migrations via `golang-migrate/migrate`, run at deploy time (migrate job/container).
- Configuration: `.env`-driven Compose configuration (DB URL, JWT signing config, hostnames/domains, rate-limit settings).
- Observability: track TTFG and connection success/disconnect patterns; structured JSON logging hooks; request/trace ID propagation where applicable.
- Rate limiting: per-IP + per-account limits (especially for lobby create/join).
- Responsive UX: mobile-first responsive layouts; touch + mouse input; recent desktop + mobile browsers.
- UX “always a way forward”: clear, friendly error messages (illegal move, full lobby, etc.) and disconnect flow that fails forward (e.g., “session lost” + CTA to join another lobby).
- Accessibility baseline: WCAG 2.1 AA basics (contrast, visible focus, keyboard-operable dialogs/tables, ARIA labels, touch targets ≥44px) plus basic a11y testing (keyboard-only + screen reader smoke + automated checks).

### FR Coverage Map

FR1: Epic 1 - Account signup (no payment)
FR2: Epic 1 - Sign in on desktop/mobile
FR3: Epic 1 - Basic profile management
FR4: Epic 2 - View available games (Connect4 now; others later)
FR5: Epic 2 - Game selection into join/host flow
FR6: Epic 2 - Browse open lobbies with seat availability
FR7: Epic 2 - Auto-join an open lobby and take a seat
FR8: Epic 2 - Create a new lobby when none available
FR9: Epic 3 - Host starts/ends a session
FR10: Epic 3 - Host removes a player (kick)
FR11: Epic 3 (Connect4) / Epic 4 (Draughts) - Rules-enforced moves; reject illegal moves
FR12: Epic 3 (Connect4) / Epic 4 (Draughts) - Authoritative shared state updates to all players
FR13: Epic 3 (Connect4) / Epic 4 (Draughts) - Game results (win/lose/draw per game rules)
FR14: Epic 2 - Copy invite link to lobby/session
FR15: Epic 2 - Invite link lands user in correct context after signup/login
FR16: Epic 5 (Deferred) - Disconnect messaging and fail-forward “rejoin new lobby” flow
FR17: Epic 6 (Deferred) - Record session events for troubleshooting/fairness
FR18: Epic 2 - Responsive, usable experience on recent desktop/mobile browsers
FR19: Epic 2 - Shareable URLs for lobby/session; deep link into context
FR20: Epic 7 (Deferred) - Track 2–4 minute onboarding funnel/metrics
FR21: Epic 3 - Connect4 move legality and turn enforcement
FR22: Epic 3 - Connect4 win/draw detection
FR23: Epic 4 - 10x10 draughts/checkers core rules + illegal move rejection
FR24: Epic 4 - Forced captures (including backward captures for men)
FR25: Epic 4 - Multi-jump captures in a single turn
FR26: Epic 4 - Kinging + flying kings
FR27: Epic 4 - Draughts outcomes incl. mutual-agreement draw
FR28: Epic 3 (Connect4) / Epic 4 (Draughts) - Clear game feedback, turn/last move/result UI

## Epic List

### Epic 1: Account Setup & Basic Profile
Users can sign up, sign in, and manage a minimal profile so they can participate in sessions.
**FRs covered:** FR1, FR2, FR3

### Epic 2: Connect4 Discovery, Lobby Join/Create, and Invite Links
Users can find Connect4, join an existing lobby or create a new one, and invite others via shareable links/URLs (desktop/mobile-friendly).
**FRs covered:** FR4, FR5, FR6, FR7, FR8, FR14, FR15, FR18, FR19

### Epic 3: Connect4 Session Control + Real-Time Gameplay
Players can start a Connect4 session, make rules-enforced moves via real-time updates, and reach a clear outcome; hosts can manage the session.
**FRs covered:** FR9, FR10, FR11, FR12, FR13, FR21, FR22, FR28

### Epic 4: Draughts 10x10 (Flying Kings) Gameplay
Users can play draughts/checkers on a 10x10 board with your selected rules (forced captures, multi-jump, backward captures for men, flying kings, mutual-agreement draw).
**FRs covered:** FR11, FR12, FR13, FR23, FR24, FR25, FR26, FR27, FR28

### Epic 5: Resilience UX — Disconnect Handling (Deferred)
If a player disconnects, the product communicates what happened and provides a “fast path” forward (MVP says no state restore).
**FRs covered:** FR16

### Epic 6: Observability & Fair Play Signals (Deferred)
The system records session events for troubleshooting and fairness checks (and supports basic operational visibility).
**FRs covered:** FR17

### Epic 7: Onboarding Speed Metrics (Deferred)
The product tracks and optimizes the “landing → seated in lobby” funnel against the target.
**FRs covered:** FR20

## Epic 1: Account Setup & Basic Profile

Users can sign up, sign in, and manage a minimal profile so they can participate in sessions.

### Story 1.1: Set up initial project from starter template

As a developer,
I want to initialize the project using the selected starter template and baseline repo structure,
So that subsequent stories can be implemented efficiently and consistently.

**Acceptance Criteria:**

**FRs:** FR1, FR2, FR3

**Given** the architecture specifies a Vite + React + TypeScript starter for the web client
**When** I initialize the web project
**Then** the repository contains a working Vite React TypeScript app (build and dev scripts runnable)
**And** the web app can be served as static assets (`web/dist`) in production

**Given** the project repository is initialized
**When** I run the basic build/test checks configured for the repo
**Then** they pass

### Story 1.2: Sign up with email + password (no payments)

As a visitor,
I want to create an account with email + password,
So that I can join lobbies and play games.

**Acceptance Criteria:**

**FRs:** FR1

**Given** I am not authenticated
**When** I submit a valid email + password to `POST /api/v1/auth/signup`
**Then** a new user account is created and the password is stored as a bcrypt hash (never plaintext)
**And** the response returns an access token (JWT) for the new user

**Given** I submit an email that is already registered
**When** I call `POST /api/v1/auth/signup`
**Then** the request is rejected with a clear error response
**And** no duplicate account is created

**Given** I submit invalid input (e.g., malformed email or too-short password)
**When** I call `POST /api/v1/auth/signup`
**Then** the request is rejected with validation errors

**Given** the signup endpoint is invoked
**When** the service logs the request
**Then** it does not log passwords (or password hashes)
**And** it does not include sensitive secrets in logs

### Story 1.3: Sign in on desktop/mobile (JWT)

As a registered user,
I want to sign in with my email + password,
So that I can access lobbies and gameplay as an authenticated user.

**Acceptance Criteria:**

**FRs:** FR2

**Given** I have a registered account
**When** I submit valid credentials to `POST /api/v1/auth/login`
**Then** I receive an access token (JWT)
**And** I can use that token to authenticate subsequent API requests

**Given** I submit invalid credentials
**When** I call `POST /api/v1/auth/login`
**Then** the request is rejected with a clear error response
**And** the response does not reveal whether the email exists

**Given** the login endpoint is invoked
**When** the service logs the request
**Then** it does not log passwords (or password hashes)
**And** it does not include sensitive secrets in logs

### Story 1.4: View/edit basic profile (username, optional avatar)

As a signed-in user,
I want to view and edit my basic profile (username, optional avatar),
So that my identity is visible and consistent in lobbies and games.

**Acceptance Criteria:**

**FRs:** FR3

**Given** I am authenticated with a valid JWT
**When** I request my profile via `GET /api/v1/me`
**Then** I receive my current profile fields (at least username and avatar if set)

**Given** I am authenticated with a valid JWT
**When** I update my profile via `PATCH /api/v1/me` with a new username (and optionally avatar)
**Then** my profile is updated and persisted
**And** a subsequent `GET /api/v1/me` reflects the changes

**Given** I submit invalid profile data (e.g., empty username)
**When** I call `PATCH /api/v1/me`
**Then** the request is rejected with validation errors

## Epic 2: Connect4 Discovery, Lobby Join/Create, and Invite Links

Users can find Connect4, join an existing lobby or create a new one, and invite others via shareable links/URLs (desktop/mobile-friendly).

### Story 2.1: List games (authenticated) with “coming soon”

As a signed-in user,
I want to view the available games on the platform,
So that I can choose what to play next.

**Acceptance Criteria:**

**FRs:** FR4

**Given** I am authenticated with a valid JWT
**When** I request the games list via `GET /api/v1/games`
**Then** I receive a list that includes at least:
**And** it shows **Connect4** as available/playable
**And** it shows **Checkers** as “coming soon” (not playable in MVP)

**Given** I am not authenticated
**When** I call `GET /api/v1/games`
**Then** the request is rejected as unauthorized

**Given** the games list endpoint is invoked
**When** the service logs the request
**Then** it does not log sensitive auth material (e.g., raw JWTs)

### Story 2.2: List public Connect4 lobbies (authenticated) with seat availability

As a signed-in user,
I want to view the list of public Connect4 lobbies and see if seats are available,
So that I can quickly join an open game.

**Acceptance Criteria:**

**FRs:** FR6

**Given** I am authenticated with a valid JWT
**When** I request the lobby list via `GET /api/v1/lobbies?game=connect4`
**Then** I receive a list of public lobbies for Connect4
**And** each lobby includes enough information to determine seat availability (e.g., max seats = 2, current occupied seats)

**Given** I am not authenticated
**When** I call `GET /api/v1/lobbies?game=connect4`
**Then** the request is rejected as unauthorized

**Given** there are no open lobbies
**When** I call `GET /api/v1/lobbies?game=connect4`
**Then** I receive an empty list (not an error)

### Story 2.3: Create a public Connect4 lobby (2 seats)

As a signed-in user,
I want to create a public Connect4 lobby with exactly 2 player seats,
So that I can invite someone or wait for another player to join.

**Acceptance Criteria:**

**FRs:** FR8

**Given** I am authenticated with a valid JWT
**When** I create a lobby via `POST /api/v1/lobbies` with `game=connect4`
**Then** a new lobby is created as public/listed
**And** the lobby’s max seats is exactly 2
**And** the response returns the created lobby (including its `lobby_id` and shareable URL)

**Given** I am not authenticated
**When** I call `POST /api/v1/lobbies`
**Then** the request is rejected as unauthorized

**Given** I attempt to create a Connect4 lobby with a seat count other than 2
**When** I call `POST /api/v1/lobbies`
**Then** the request is rejected with a validation error

### Story 2.4: Join a Connect4 lobby and claim a seat

As a signed-in user,
I want to join a specific Connect4 lobby and claim an available seat,
So that I can participate in the upcoming session.

**Acceptance Criteria:**

**FRs:** FR7

**Given** I am authenticated with a valid JWT
**When** I join a lobby via `POST /api/v1/lobbies/{lobby_id}/join`
**Then** I become a lobby member assigned to an available player seat

**Given** the lobby is full (both seats taken)
**When** I call `POST /api/v1/lobbies/{lobby_id}/join`
**Then** the request is rejected with a clear “lobby full” error

**Given** I already joined the lobby previously
**When** I call `POST /api/v1/lobbies/{lobby_id}/join` again
**Then** the operation is idempotent (no duplicate membership)
**And** I remain assigned to my existing seat

**Given** I am not authenticated
**When** I call `POST /api/v1/lobbies/{lobby_id}/join`
**Then** the request is rejected as unauthorized

### Story 2.5: Auto-join an open Connect4 lobby (or create one)

As a signed-in user,
I want to automatically get seated in an open Connect4 lobby if one exists,
So that I can start playing with minimal friction.

**Acceptance Criteria:**

**FRs:** FR7, FR8

**Given** I am authenticated with a valid JWT
**When** I request auto-join via `POST /api/v1/lobbies/auto-join` with `game=connect4`
**Then** if an open lobby with an available seat exists, I am joined to that lobby and assigned a seat
**And** the response includes the selected `lobby_id` and shareable lobby URL

**Given** there is no open lobby with available seats
**When** I call `POST /api/v1/lobbies/auto-join` with `game=connect4`
**Then** a new public Connect4 lobby is created (2 seats)
**And** I am joined to it and assigned a seat
**And** the response returns the created `lobby_id` and shareable lobby URL

**Given** I am not authenticated
**When** I call `POST /api/v1/lobbies/auto-join`
**Then** the request is rejected as unauthorized

### Story 2.6: Lobby invite link + deep-link after login

As a signed-in user,
I want a shareable lobby invite link that opens the correct lobby after login,
So that I can bring another player into my Connect4 lobby.

**Acceptance Criteria:**

**FRs:** FR14, FR15, FR19

**Given** I am authenticated and viewing a lobby
**When** I request the lobby details via `GET /api/v1/lobbies/{lobby_id}`
**Then** the response includes a stable, shareable lobby URL for that lobby

**Given** I am not authenticated
**When** I navigate to a lobby URL (deep link)
**Then** I am required to log in
**And** after successful login, I am redirected back into the same lobby context

**Given** the lobby does not exist (or is invalid)
**When** I navigate to its lobby URL
**Then** I see a clear “not found” experience with a path back to game/lobby selection

### Story 2.7: Responsive UI for games + lobby list + join/create + copy invite

As a signed-in user,
I want a responsive UI to select a game and join/create a Connect4 lobby (and copy an invite link),
So that I can quickly get into a lobby on desktop or mobile.

**Acceptance Criteria:**

**FRs:** FR4, FR5, FR6, FR7, FR8, FR14, FR18, FR19

**Given** I am authenticated
**When** I open the Games screen
**Then** I see Connect4 as available
**And** I see Checkers labeled “coming soon”

**Given** I select Connect4
**When** I view the lobby list
**Then** I see the list of public Connect4 lobbies with seat availability
**And** I can join an open lobby or create a new lobby

**Given** I am in a lobby
**When** I click “Copy invite link”
**Then** the lobby URL is copied
**And** I get a clear confirmation message

**Given** I use a mobile browser
**When** I view the lobby list UI
**Then** it remains usable and readable (responsive layout)

## Epic 3: Connect4 Session Control + Real-Time Gameplay

Players can start a Connect4 session, make rules-enforced moves via real-time updates, and reach a clear outcome; hosts can manage the session.

### Story 3.1: Host starts a Connect4 session from a lobby (random Player 1)

As a lobby host,
I want to start a Connect4 session once both player seats are filled,
So that the game can begin with a clear first turn assignment.

**Acceptance Criteria:**

**FRs:** FR9

**Given** I am authenticated with a valid JWT
**And** I am the host of a Connect4 lobby
**And** the lobby has exactly 2 players seated
**When** I click Start and the client calls `POST /api/v1/lobbies/{lobby_id}/start`
**Then** a new Connect4 session is created for that lobby
**And** Player 1 is assigned randomly between the two seated players
**And** the initial game state is created and stored as the authoritative state
**And** the response includes the created `session_id`

**Given** I am not the lobby host
**When** I call `POST /api/v1/lobbies/{lobby_id}/start`
**Then** the request is rejected as forbidden

**Given** the lobby does not have 2 seated players
**When** I call `POST /api/v1/lobbies/{lobby_id}/start`
**Then** the request is rejected with a clear error response (e.g., “not enough players”)

### Story 3.2: Connect4 session WebSocket connect + auth + state subscription

As a signed-in player,
I want to connect to the session WebSocket and receive authoritative game state updates,
So that my client stays in sync during real-time play.

**Acceptance Criteria:**

**FRs:** FR12

**Given** I am authenticated with a valid JWT
**When** my client connects to the session WebSocket using `?access_token=<jwt>` and identifies the `session_id`
**Then** the connection is accepted and I begin receiving server-sent messages for that session

**Given** I provide an invalid/expired token
**When** I attempt to connect
**Then** the connection is rejected as unauthorized

**Given** the server sends any WebSocket message
**When** the client receives it
**Then** it follows the mandatory message envelope (e.g., includes `type`, `payload`, `ts`, `seq`)
**And** the session’s canonical state message includes a `schema_version`

**Given** the WebSocket connection is established
**When** the service logs the connection attempt
**Then** it must not log the full URL/query string (to avoid token leakage)

### Story 3.3: Make a Connect4 move (validate + broadcast state)

As a signed-in player in an active Connect4 session,
I want to submit a move by choosing a column,
So that the server applies legal moves, rejects illegal ones, and keeps all players in sync.

**Acceptance Criteria:**

**FRs:** FR11, FR12, FR21

**Given** I am authenticated with a valid JWT
**And** I am a seated player in the session
**And** it is my turn
**When** I submit a move (a column index) via the session WebSocket
**Then** the server validates the move (column exists, column not full, correct turn)
**And** applies the move to the authoritative state
**And** broadcasts the updated canonical game state to all session players

**Given** it is not my turn, or the column is invalid/full
**When** I submit a move
**Then** the move is rejected with a clear, game-specific error message
**And** the authoritative state is not changed

**Given** the server broadcasts a post-move state update
**When** clients receive it
**Then** it includes the updated board, current turn, and last move
**And** it increments `seq` monotonically for that session stream

### Story 3.4: Detect Connect4 win/draw and end the game

As a player in a Connect4 session,
I want the system to detect a win or draw and clearly end the game,
So that the result is unambiguous for both players.

**Acceptance Criteria:**

**FRs:** FR13, FR22

**Given** a valid move is applied
**When** the updated board results in 4-in-a-row for the current player (horizontal/vertical/diagonal)
**Then** the server marks the session as ended with that player as the winner
**And** broadcasts a result state to all players

**Given** a valid move is applied
**When** the board becomes full and no 4-in-a-row exists
**Then** the server marks the session as ended as a draw
**And** broadcasts a draw result state to all players

**Given** the game is ended
**When** any player attempts to submit another move
**Then** the move is rejected with a clear error
**And** the authoritative state is not changed

### Story 3.5: Host ends the session at any time

As a session host,
I want to end a Connect4 session at any time,
So that I can stop the game if needed.

**Acceptance Criteria:**

**FRs:** FR9

**Given** I am authenticated with a valid JWT
**And** I am the host for the session’s lobby
**When** I request session end (API or WebSocket control message)
**Then** the session is marked ended
**And** all connected players receive a “session ended” message/state update

**Given** I am not the host
**When** I attempt to end the session
**Then** the request is rejected as forbidden
**And** the session continues

### Story 3.6: Host kicks a player at any time

As a session host,
I want to kick a player from the lobby/session at any time,
So that I can remove disruptive or idle players.

**Acceptance Criteria:**

**FRs:** FR10

**Given** I am authenticated with a valid JWT
**And** I am the host of the lobby/session
**When** I kick a player (API or WebSocket control message specifying the target player)
**Then** the target player is removed from the lobby/session
**And** the kicked player receives a clear notification they were removed
**And** the remaining player(s) receive an updated lobby/session state

**Given** I am not the host
**When** I attempt to kick a player
**Then** the request is rejected as forbidden

**Given** I kick a player during active play
**When** the kick is applied
**Then** the session is ended (since Connect4 requires 2 players)
**And** the remaining player sees a clear “opponent removed / session ended” state

### Story 3.7: Gameplay UI (board + turn + errors + result)

As a signed-in player,
I want a clear Connect4 gameplay UI that stays in sync with the server,
So that I can play moves, understand turns, and see results.

**Acceptance Criteria:**

**FRs:** FR28, FR18

**Given** I am in an active Connect4 session
**When** the server broadcasts the canonical game state
**Then** the UI renders the board correctly
**And** shows whose turn it is
**And** shows the last move

**Given** it is my turn
**When** I select a valid column
**Then** my client sends the move to the server
**And** the UI updates based on the server’s subsequent state update (not speculative-only)

**Given** I attempt an illegal move (e.g., full column or not my turn)
**When** the server rejects the move
**Then** I see a clear, game-specific error message
**And** the board state does not incorrectly change

**Given** the game ends (win/draw/session ended)
**When** the result is broadcast
**Then** the UI clearly shows the final result state

## Epic 4: Draughts 10x10 (Flying Kings) Gameplay

Users can play draughts/checkers on a 10x10 board with your selected rules (forced captures, maximum-capture priority, multi-jump, backward captures for men, flying kings, mutual-agreement draw).

### Story 4.1: Initialize 10x10 draughts state (standard setup, white starts)

As a signed-in player,
I want the system to initialize a correct 10x10 draughts game state,
So that a match can start with a valid board and turn order.

**Acceptance Criteria:**

**FRs:** FR23

**Given** a new draughts session is created
**When** the server initializes the game state
**Then** the board is 10x10 using playable dark squares
**And** each side starts with 20 men placed in the standard 10x10 setup (first 4 rows per side on playable squares)
**And** no kings are present at start
**And** it is White’s turn first
**And** the canonical state includes a `schema_version`

**Given** a draughts session has an initialized state
**When** the server publishes the initial canonical state to clients
**Then** the client can render the correct piece counts and positions
**And** the UI can display that it is White’s turn

### Story 4.2: Validate non-capture moves (men forward-only; flying kings)

As a signed-in player in a draughts session,
I want the server to validate non-capture moves according to the rules,
So that only legal moves are applied.

**Acceptance Criteria:**

**FRs:** FR23

**Given** it is my turn and no capture is available for my pieces
**When** I submit a non-capture move
**Then** the server accepts it only if it is legal:
**And** men move diagonally forward by one square to an empty playable square
**And** kings can move (“fly”) diagonally any number of empty squares
**And** the move stays on playable dark squares

**Given** it is my turn and I attempt an illegal non-capture move
**When** I submit the move
**Then** the server rejects it with a clear, game-specific error
**And** the authoritative state is not changed

### Story 4.3: Validate captures + enforce forced capture (including backward men captures)

As a signed-in player in a draughts session,
I want the server to validate capture moves and enforce forced captures,
So that captures follow the rules and players cannot ignore them.

**Acceptance Criteria:**

**FRs:** FR24

**Given** it is my turn and at least one capture is available for my pieces
**When** I attempt a non-capture move
**Then** the server rejects it with a clear “capture required” error
**And** the authoritative state is not changed

**Given** I submit a capture move with a man
**When** the capture is validated
**Then** the server accepts only diagonal jump captures over an adjacent opponent piece to an empty landing square
**And** men may capture backward as well as forward

**Given** I submit a capture move with a flying king
**When** the capture is validated
**Then** the king may capture by jumping over exactly one opponent piece on a diagonal
**And** may land on any empty square beyond the captured piece on that diagonal
**And** no other pieces may be jumped in that same jump segment

### Story 4.4: Enforce maximum-capture rule

As a signed-in player in a draughts session,
I want the server to require that I choose a capture sequence that captures the maximum number of pieces,
So that the maximum-capture rule is enforced.

**Acceptance Criteria:**

**FRs:** FR24

**Given** it is my turn and multiple capture sequences are available
**When** I submit a capture move/sequence that does not result in the maximum capture count
**Then** the server rejects it with a clear “maximum capture required” error
**And** the authoritative state is not changed

**Given** it is my turn and I submit a capture move/sequence that achieves the maximum capture count
**When** the move/sequence is validated
**Then** the server accepts it and applies it to the authoritative state

### Story 4.5: Multi-jump capture continuation + kinging timing

As a signed-in player in a draughts session,
I want the server to enforce multi-jump captures and apply kinging correctly,
So that capture turns are completed properly.

**Acceptance Criteria:**

**FRs:** FR25, FR26

**Given** I make a capture and another capture is available for that same piece
**When** the capture is applied
**Then** the server requires that I continue capturing with that same piece
**And** rejects attempts to end the turn early

**Given** my capturing piece lands on the king row during a capture sequence
**When** the move sequence ends
**Then** the piece becomes a king only if it stops on the king row at the end of the turn
**And** it moves as a king starting on the next turn (not mid-turn)

### Story 4.6: Draughts outcomes + mutual-agreement draw

As a signed-in player in a draughts session,
I want the game to end with a clear outcome (win/loss/draw),
So that the result is unambiguous.

**Acceptance Criteria:**

**FRs:** FR27

**Given** my opponent has no legal moves (or no pieces)
**When** the server evaluates the game state
**Then** the session ends and I am declared the winner
**And** all players receive the result state

**Given** I propose a draw during an active game
**When** the other player accepts the draw
**Then** the session ends as a draw
**And** all players receive the draw result state

**Given** I propose a draw and the other player declines (or does nothing)
**When** play continues
**Then** the game does not end as a draw

### Story 4.7: Draughts UI (click-to-move + optional legal-move highlighting)

As a signed-in player,
I want a draughts UI that supports click-to-move and (optionally) highlights legal moves,
So that I can play the 10x10 ruleset correctly and understand errors/results.

**Acceptance Criteria:**

**FRs:** FR28

**Given** I am in an active draughts session
**When** the server broadcasts the canonical state
**Then** the UI renders the 10x10 board and piece positions correctly
**And** shows whose turn it is

**Given** I select one of my pieces
**When** legal-move highlighting is enabled
**Then** the UI displays the legal destinations for that piece (respecting forced capture + maximum-capture rules)

**Given** legal-move highlighting is disabled
**When** I select one of my pieces
**Then** the UI still allows click-to-move by selecting a destination square
**And** illegal attempts are handled via server error messages

**Given** I attempt an illegal move (wrong turn / capture required / maximum capture required / invalid landing)
**When** the server rejects the move
**Then** I see a clear, game-specific error message
**And** the board state does not incorrectly change

**Given** the game ends (win/loss/draw)
**When** the result is broadcast
**Then** the UI clearly shows the final result state

## Epic 5: Resilience UX — Disconnect Handling (Deferred)

If a player disconnects, the system communicates loss of session and offers a fast path to continue elsewhere (no state restore in MVP).

### Story 5.1: Detect disconnect and end session (no winner)

As a player in an active session,
I want the system to handle an opponent disconnect clearly,
So that the session ends cleanly with no confusion.

**Acceptance Criteria:**

**FRs:** FR16

**Given** a Connect4 session is active with 2 connected players
**When** one player disconnects (WebSocket connection closes)
**Then** the server marks the session as ended with “no winner”
**And** broadcasts a “player disconnected / session ended” state to the remaining player

**Given** the disconnected player reconnects later
**When** they attempt to rejoin
**Then** they are not restored into the ended session
**And** they are guided to start/join a new lobby instead (no state restore in MVP)

### Story 5.2: Client UX for disconnect (message + CTA to Games)

As a player,
I want a clear UX when I disconnect or lose my session,
So that I can quickly get back to a working flow.

**Acceptance Criteria:**

**FRs:** FR16

**Given** I am disconnected from the session (or I receive a “session ended due to disconnect” state)
**When** the client detects this condition
**Then** I see a clear message explaining the session was lost/ended
**And** I see a primary CTA that takes me to the Games screen

**Given** I click the CTA
**When** navigation occurs
**Then** I land on the Games screen and can start/join a new lobby

### Story 5.3: Lobby disconnect handling (pre-start seat cleanup)

As a lobby participant,
I want the lobby to stay accurate when someone disconnects before the game starts,
So that seats free up and others can still play.

**Acceptance Criteria:**

**FRs:** FR16

**Given** a Connect4 lobby exists and the session has not started
**When** a seated player disconnects/leaves
**Then** the player is removed from the lobby membership (seat is freed)
**And** the lobby list reflects updated seat availability

**Given** the host disconnects before start
**When** the disconnect is handled
**Then** the lobby is closed
**And** any remaining player is informed and returned to the Games screen

## Epic 6: Observability & Fair Play Signals (Deferred)

The system records session events (join, start, moves, end, kicks) for basic troubleshooting and fairness checks.

### Story 6.1: Persist session events in Postgres

As a platform operator,
I want session events to be recorded in Postgres,
So that we have an audit trail for troubleshooting and basic fairness checks.

**Acceptance Criteria:**

**FRs:** FR17

**Given** the system persists session events
**When** a relevant event occurs (join/start/move/end/kick)
**Then** a row is inserted into a `session_events` table with at least:
**And** `event_id`, `session_id`, `event_type`, `actor_user_id` (nullable when not applicable), `payload` (JSON), `created_at`

**Given** events are recorded over time
**When** querying by `session_id` ordered by `created_at` (and/or sequence)
**Then** the full event history can be reconstructed for that session

### Story 6.2: Record join/start/end/kick events

As a platform operator,
I want the system to record key lobby/session lifecycle events,
So that we can diagnose issues and review basic session timelines.

**Acceptance Criteria:**

**FRs:** FR17

**Given** a player joins a lobby (claims a seat)
**When** the join is accepted
**Then** a `session_events` row is recorded with `event_type=join` and relevant identifiers in `payload`

**Given** the host starts a session
**When** the session is created
**Then** a `session_events` row is recorded with `event_type=start`

**Given** a session ends (normal end, host end, disconnect end)
**When** the session is marked ended
**Then** a `session_events` row is recorded with `event_type=end` including an end reason in `payload`

**Given** the host kicks a player
**When** the kick is applied
**Then** a `session_events` row is recorded with `event_type=kick` including the target user id in `payload`

### Story 6.3: Record move events with minimal payload + ordering

As a platform operator,
I want each accepted move to be recorded as an event,
So that we can reconstruct gameplay sequences for troubleshooting.

**Acceptance Criteria:**

**FRs:** FR17

**Given** a Connect4 move is accepted
**When** the authoritative state is updated
**Then** a `session_events` row is recorded with `event_type=move`
**And** `payload` contains at least the move input (e.g., column) and resulting `seq`

**Given** a draughts move/capture sequence is accepted
**When** the authoritative state is updated
**Then** a `session_events` row is recorded with `event_type=move`
**And** `payload` contains at least the move representation (from/to, captures, promotion if any) and resulting `seq`

**Given** events exist for a session
**When** ordered by `created_at` and/or `seq`
**Then** the move timeline is reconstructable without ambiguity

## Epic 7: Onboarding Speed Metrics (Deferred)

The product measures onboarding speed from login completion to session start, using anonymized metrics in Postgres.

### Story 7.1: Persist anonymized onboarding funnel events

As a platform operator,
I want onboarding funnel events recorded in Postgres without user-identifying data,
So that we can measure signup-to-play speed while respecting privacy.

**Acceptance Criteria:**

**FRs:** FR20

**Given** a user completes login successfully
**When** the system processes the login success
**Then** an onboarding metric event is stored with an anonymized attempt identifier (no `user_id`)
**And** the event includes a timestamp and `event_type=login_success`

**Given** a session is started successfully
**When** the host starts the session
**Then** an onboarding metric event is stored with the same anonymized attempt identifier
**And** the event includes a timestamp and `event_type=session_started`

### Story 7.2: Compute and store onboarding duration + basic aggregates

As a platform operator,
I want to compute the time from login success to session start per attempt,
So that we can track whether onboarding meets the target.

**Acceptance Criteria:**

**FRs:** FR20

**Given** an anonymized attempt has both `login_success` and `session_started` events
**When** the second event is recorded (or via a background aggregation process)
**Then** the system stores the computed duration for that attempt in Postgres

**Given** durations are stored over time
**When** we query aggregates
**Then** we can compute basic statistics (e.g., p50 and p95) for the login→session-start funnel

### Story 7.3: Record metrics for all main entry paths

As a platform operator,
I want onboarding metrics to be recorded consistently across entry paths,
So that our funnel measurement isn’t biased.

**Acceptance Criteria:**

**FRs:** FR20

**Given** a user reaches session start via manual lobby join
**When** they complete login and the session starts
**Then** the same `login_success` → `session_started` attempt is recorded

**Given** a user reaches session start via auto-join
**When** they complete login and the session starts
**Then** the same `login_success` → `session_started` attempt is recorded

**Given** a user reaches session start via a lobby deep link (invite URL)
**When** they complete login and the session starts
**Then** the same `login_success` → `session_started` attempt is recorded
