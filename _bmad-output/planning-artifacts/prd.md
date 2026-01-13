---
stepsCompleted: ['step-01-init', 'step-02-discovery', 'step-03-success', 'step-04-journeys', 'step-05-domain', 'step-06-innovation', 'step-07-project-type', 'step-08-scoping', 'step-09-functional', 'step-10-nonfunctional', 'step-11-polish', 'step-12-complete']
inputDocuments:
  - _bmad-output/planning-artifacts/product-brief-mygameplatform-2026-01-11.md
  - _bmad-output/planning-artifacts/research/market-multiplayer-gaming-platform-market-research-2026-01-11.md
  - _bmad-output/analysis/brainstorming-session-2026-01-10_22-24-37.md
  - _bmad-output/analysis/brainstorming-session-2026-01-10_23-59-43.md
  - _bmad-output/analysis/brainstorming-session-2026-01-11.md
workflowType: 'prd'
classification:
  projectType: web_app
  domain: general
  complexity: low
  projectContext: greenfield
workflow_completed: true
workflow: 'edit'
lastEdited: '2026-01-13'
editHistory:
  - date: '2026-01-13'
    changes: 'Replaced Checkers with Draughts (10x10) worldwide canonical ruleset; added explicit FR21–FR28 for Connect4 + International Draughts rules/outcomes; resolved template placeholders in title/author/date.'
date: '2026-01-13'
---

# Product Requirements Document - mygameplatform

**Author:** Laurent
**Date:** 2026-01-13

## Executive Summary

Forever-free multiplayer board-game platform for casual players: create an account fast, join or host an active lobby, and play rules-enforced games on desktop/mobile browsers via shareable URLs.

MVP games: Connect4 and **Draughts (10x10)**. Draughts uses a single canonical worldwide ruleset: International Draughts (forced captures, maximum-capture priority, multi-jumps, flying kings, draw by mutual agreement).

Differentiator: no paywalls or monetization friction; optimize time-to-first-game and lobby liquidity for drop-in sessions.

## Project Classification
- Project type: web_app (MPA with real-time play)
- Domain: general; complexity: low; context: greenfield

## Success Criteria

### User Success
- Time-to-first-game (TTFG) p95 ≤ 4 minutes from landing to game start.
- Users can complete a full game session without forced disconnects; basic rule validation ensures fair play.
- Weekly return rate tracked for core casuals (baseline to be measured in MVP).

### Business Success
- Reach WAU 100 and MAU 200 within MVP window.
- Maintain active lobbies during peak casual hours (liquidity target to be baselined).
- Invite funnel: invites sent/accepted tracked; invite acceptance rate baselined in MVP.

### Technical Success
- Baseline availability (no strict SLO) with reconnect success rate measured; monitor disconnect-induced churn.
- Rule-validation correctness for shipped games (Connect4, Draughts (10x10)).
- Page-to-WebSocket connect time tracked to inform latency improvements.

### Measurable Outcomes
- TTFG p95 ≤ 4 minutes; median tracked.
- WAU 100 / MAU 200.
- Sessions per active user/week (baseline & trend).
- Invite acceptance rate (baseline & trend).
- Reconnect success % (tracked, no target yet).

## Product Scope

### MVP - Minimum Viable Product
- Account creation and instant play (2–4 minute signup; join/host quickly).
- Games: Connect4, Draughts (10x10).
- Lobby/session flow: browse/auto-join open lobbies; host start/end; basic rule validation.
- Cross-device: responsive web (desktop + mobile browser).

### Growth Features (Post-MVP)
- Reconnect grace window and transport health hooks.
- Stronger moderation/abuse telemetry and report flow.
- Additional games beyond Connect4/Draughts (10x10)

### Vision (Future)
- Broader catalog with day-and-date digital launches.
- Richer social features (friends, events, ladders) and reconnect UX.
- Optional monetization/cosmetics if “forever free” stance is revisited.

## User Journeys

### Journey 1: Casual Player — Happy Path (Drop-in)
- Device: mobile browser; finds via friend link.
- Flow: Opens link → 2-minute signup → sees open lobbies → auto-joins → plays Connect4 to completion (<1h session).
- Aha: Near-immediate join with no paywall; rules enforced; finishes a game.
- Outcome: Returns weekly; shares invite link to a friend next time.

### Journey 2: Casual Player — Edge (Disconnect without reconnect feature)
- Device: mobile browser; joins an open lobby.
- Event: Connection drops mid-game; reconnect not available in MVP.
- Recovery: User relaunches, re-joins a new lobby quickly; learns progress isn’t saved.
- Requirement signal: Need clear messaging on loss, fast rejoin path, protect experience by making new join fast; future reconnect is a growth item.

### Journey 3: Host — Quick Table Setup for Friends
- Device: desktop browser.
- Flow: Logs in → creates lobby for Connect4/Draughts (10x10) → shares invite link → starts when two seats filled → can end session.
- Aha: Can spin up a game and start within minutes; no payments or complex settings.
- Outcome: Uses platform for casual friend/family sessions weekly.

### Journey 4 (Deferred for later): Indie Publisher
- Not in MVP flow; to be detailed post-MVP when adding more games.

### Journey Requirements Summary
- Instant signup + lobby discovery/auto-join.
- Host controls: create lobby, start/end, share invites.
- Clear handling of disconnects in MVP (no reconnect): messaging and fast rejoin path.
- Rule validation for Connect4/Draughts (10x10).
- Cross-device responsive web (mobile+desktop).

## Web App Specific Requirements

### Project-Type Overview
- Web app (MPA) with per-session URLs for shareable lobbies/games.
- Real-time via WebSockets.

### Technical Architecture Considerations
- MPA with distinct routes/pages (lobby list, lobby detail/session).
- Recent desktop/mobile browsers (evergreen) as baseline support.
- Minimal SEO; focus on shareable URLs and invite landing performance.
- Real-time transport: WebSockets; page-to-WS connect time tracked.

### Browser Matrix
- Support: latest Chrome/Firefox/Edge (desktop), latest Safari (mobile/desktop), Chrome mobile.
- No requirement to support legacy/IE-class browsers.

### Responsive Design
- Comfortable on mobile portrait and desktop; touch and mouse input.

### Performance Targets
- Time-to-first-game p95 ≤ 4 minutes (aligned to success metrics); measure page-to-WS connect.
- Keep page bundles lightweight per MPA page.

### SEO Strategy
- Minimal: shareable lobby/game URLs render title/description for previews; basic indexable landing page.

### Accessibility Level
- No formal accessibility requirement for MVP; note as future consideration.

## Project Scoping & Phased Development

### MVP Strategy & Philosophy
- Approach: Experience MVP — prove instant, free play with minimal friction.
- Team (implied): small web team + light backend; real-time + basic ops.

### MVP Feature Set (Phase 1)
- Journeys: Casual happy path; Host quick setup; Edge: disconnect (messaging/fast rejoin to new lobby).
- Must-haves: 2–4 min signup-to-play; lobby list + auto-join; host start/end; rule validation for Connect4/Draughts (10x10); responsive web; WebSockets; invite links.

### Post-MVP Features (Phase 2)
- Reconnect grace + transport health hooks.
- Moderation/abuse telemetry & report flow.
- Additional games (mono or multi-player, more complex than connect4 and draughts)

### Phase 3 (Expansion)
- Richer social (friends/events/ladders).
- Creator/publisher onboarding; day-and-date launches.
- Optional monetization/cosmetics if strategy shifts.

### Risk Mitigation Strategy
- Technical: start with two simple games; measure page-to-WS connect and TTFG; avoid reconnect in MVP to de-risk.
- Market: validate WAU/MAU and invite acceptance; ensure lobby liquidity with fast join + minimal game set.
- Resource: keep MPA pages light; defer chat/payments/moderation depth to Phase 2.

## Functional Requirements

### Account & Identity
- FR1: A visitor can create an account in one flow without providing payment information.
- FR2: A registered user can sign in on desktop or mobile browser.
- FR3: A user can view and edit basic profile fields (username, avatar optional) without payment details.

### Game Catalog & Selection
- FR4: A user can view available games (Connect4, Draughts (10x10)) with basic rules/players info.
- FR5: A user can pick a game and proceed directly to join/host flow.

### Lobby Discovery & Joining
- FR6: A user can see a list of open lobbies for a selected game with seat availability.
- FR7: A user can auto-join an open lobby and take an available seat.
- FR8: If no suitable lobby exists, a user can create a new lobby for that game.

### Hosting & Session Control
- FR9: A host can start a session once minimum seats are filled and end the session at any time.
- FR10: A host can remove a player from the session before or during play.

### Gameplay (Authoritative, Rules-Enforced)
- FR11: Players in a session can make moves according to game rules; illegal moves are rejected.
- FR12: The system updates and shares canonical game state to all session players after each valid move.
- FR13: A session concludes with a clear game result (win/lose/draw) per rules for Connect4 and Draughts (10x10).

#### Draughts (10x10) Ruleset (International Draughts)
- FR21: The system implements Connect4 rules (7x6 board, gravity drop into a column, alternating turns) and rejects illegal moves (invalid column, full column, wrong turn).
- FR22: The system detects and communicates Connect4 outcomes (4-in-a-row horizontal/vertical/diagonal, draw when board is full).
- FR23: The system implements draughts gameplay on a 10x10 board (International Draughts ruleset) and rejects illegal moves.
- FR24: The system enforces forced captures for draughts (a capture must be made when available), including backward captures for men.
- FR25: The system enforces maximum-capture priority for draughts when multiple capture options exist (player must choose a capture sequence that captures the maximum number of pieces).
- FR26: The system enforces multi-jump capture sequences for draughts within a single turn when available (continuation captures required until none remain).
- FR27: The system performs kinging for draughts when a piece reaches the last rank, and kings have “flying” movement/captures (multi-square diagonals) per International Draughts rules.
- FR28: The system detects and communicates draughts outcomes (win/loss), and supports draws when both players agree to a draw.

### Invites & Sharing
- FR14: A host or player can copy an invite link to a lobby/session for others to join.
- FR15: An invited user can follow a link and land directly in the correct game/lobby context after signup/login.

### Rejoin/Recovery (MVP messaging)
- FR16: If a player disconnects, the system communicates loss of session and offers fast path to rejoin a new lobby (no state restore in MVP).

### Observability & Fair Play Signals
- FR17: The system records session events (join, start, moves, end, kicks) for basic troubleshooting and fairness checks.

### Web Experience
- FR18: The product is usable on recent desktop and mobile browsers with responsive layouts and touch/mouse input.
- FR19: Each lobby/session has a shareable URL; users can load it directly to join or host if eligible.

### Onboarding Speed
- FR20: The flow from landing to sitting in a lobby enables completion within the 2–4 minute target (as tracked in metrics).

## Non-Functional Requirements

### Performance
- Time-to-first-game p95 ≤ 4 minutes (tracked); page-to-WebSocket connect time monitored.
- Game state updates visible to players within a small real-time window (WebSocket-based); no strict SLO set in MVP.

### Security
- Accounts require authenticated access; no payments collected in MVP.
- Session access via invite link must validate eligibility (seat available, game context).
- Basic protection against lobby abuse: sanitize lobby names; rate limit lobby create/join attempts.

### Scalability
- MVP targets WAU 100 / MAU 200; design to add more lobbies/games without architectural change.

### Accessibility
- No formal requirement for MVP; responsive layouts and touch/mouse input supported.

### Integration
- None required in MVP (no external payments or third-party systems).
