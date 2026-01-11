---
stepsCompleted: [1, 2, 3, 4]
inputDocuments:
  - "_bmad-output/analysis/brainstorming-session-2026-01-10_22-24-37.md"
date: 2026-01-10_23-13-35
author: Laurent
---

# Product Brief: {{project_name}}

<!-- Content will be appended sequentially through collaborative workflow steps -->

## Executive Summary

MyGamePlatform is a browser-based multiplayer gaming platform focused on turn-based games (starting with simple titles like Connect4) that lets people play together online without ads or paid tiers. The platform is intentionally designed to be truly free and self-funded, removing monetization pressure from the experience.

The MVP emphasizes a low-friction “show up and play” flow: players create accounts, browse public lobbies, and join game sessions live. Friends can meet up by finding the host’s lobby in the public list (the lobby name is the host’s name).

---

## Core Vision

### Problem Statement

Many online board-game platforms introduce ads, paywalls, or other monetization-driven friction that gets in the way of simply playing games with friends online. Even when games are available, it can be cumbersome to reliably find and join friends’ sessions.

### Problem Impact

For players who just want a fun, straightforward way to play favorite games online, ads and paywalls degrade the experience and limit access. Over time, these constraints reduce spontaneous play and make “game night” coordination harder than it should be.

### Why Existing Solutions Fall Short

Current platforms commonly rely on monetization schemes (ads, subscriptions, paid content, or upsells) that create friction and detract from a clean play experience. They may also make it harder than necessary to rendezvous with friends quickly and confidently.

### Proposed Solution

Build a truly free, browser-based multiplayer platform for turn-based games with an authoritative server that enforces rules. Players create accounts (email/password), browse public lobbies, and join sessions that auto-start when full. The host retains minimal control to end a session and replace disconnected players. Gameplay is live and synchronous, using WebSockets for real-time communication, and the game UI is responsive for desktop and mobile.

### Key Differentiators

- **Truly free:** no ads and no paid tiers; self-funded to keep incentives aligned with player enjoyment.
- **Low-friction meetup:** public lobbies with host-named lobbies makes it easy for friends to find each other.
- **Fair play by design:** authoritative server validates moves and enforces rules.
- **MVP clarity:** start with perfect-information games and a Connect4 vertical slice to validate the end-to-end experience.

## Target Users

### Primary Users

**Persona: Michel (43) — casual player**

- **Context:** Plays from both desktop and mobile; plays with a mix of existing friends and strangers.
- **Cadence:** Roughly once per week for ~30 minutes.
- **Preferences:** Simple games without complex rules (e.g., Connect4/checkers-style experiences).
- **Core pain today:** Waiting time and friction getting into a game (in addition to ads/paywalls elsewhere).
- **Onboarding tolerance:** Will register once, as long as it’s fast.
- **Success / “aha” moment:** “In a few minutes I spent a good time playing with some people.”

### Secondary Users

N/A for MVP (hardcore players are a potential later segment).

### User Journey

**Discovery**
- Finds `mygameplatform` via search or a friend recommendation.

**Onboarding**
- Creates an account quickly (email/password, minimal profile).

**Core Usage**
- Browses public lobbies and joins a session with minimal waiting time.
- Plays a simple, live synchronous, turn-based game with rules enforced by the platform.

**Success Moment**
- Gets into a game within minutes and has a fun 30-minute session with friends/others.

**Long-term**
- Returns weekly as a dependable “no-ads, no-paywalls” place to play quick games.

## Success Metrics

### User Success Metrics (MVP)

- **Time to first game started:** ≤ 3 minutes (from login to match start)
- **Lobby waiting time:** ≤ 1 minute (from join to match start)
- **Match completion rate:** ≥ 30% of started matches reach a valid end state
- **Illegal-move rejection accuracy:** 100% (authoritative referee rejects illegal moves reliably)
- **Weekly Active Users (WAU):** 100
- **Week-1 retention:** 20% (users who return the next week)

### Business Objectives

- Deliver a genuinely free experience (no ads, no paid tiers) while building a sustainable user base through trust and delight.

### Key Performance Indicators

- **Monthly Active Users (MAU):** 200
- **Operating cost:** ≤ 200 EUR/month (hosting + core ops)
- **Support burden:** ≤ 4 issues/week

## MVP Scope

### Core Features

- **Account access (minimal):** signup + login (email/password) to participate.
- **Public lobbies:** users browse a public lobby list to find active sessions.
- **Meetup mechanism:** lobby is named after the host; friends can locate the host’s lobby in the list.
- **2-player Connect4 gameplay:** live synchronous, turn-based match with an authoritative server validating moves and enforcing rules.
- **Core screens (required):** login/register, lobby list, create lobby, join lobby, game screen (Connect4), results screen (win/draw).

### MVP Scope Decisions (ADR-style)

- **ADR-001: Session discovery**
  - **Options:** invite links/codes | friend list | public lobbies
  - **Decision:** public lobbies (primary)
  - **Rationale:** lowest friction; avoids building social graph/invite mechanics early

- **ADR-002: Session privacy**
  - **Options:** public + private | public-only | private-only
  - **Decision:** public-only
  - **Rationale:** reduces complexity; optimizes “play immediately” goal

- **ADR-003: Joining behavior**
  - **Options:** choose seat | auto-assign seat | host approval
  - **Decision:** auto-assign seat
  - **Rationale:** fastest path from “join” to “game start”

- **ADR-004: Start rule**
  - **Options:** host starts | auto-start when full | configurable
  - **Decision:** auto-start when full
  - **Rationale:** minimizes waiting and coordination overhead

- **ADR-005: Authority + rules enforcement**
  - **Options:** client-authoritative | server-authoritative | hybrid
  - **Decision:** server-authoritative referee (reject illegal moves)
  - **Rationale:** fairness and correctness are non-negotiable

- **ADR-006: Real-time transport**
  - **Options:** WebSockets | WebRTC | polling
  - **Decision:** WebSockets
  - **Rationale:** good fit for live synchronous gameplay and broad browser support

- **ADR-007: Disconnect handling**
  - **Options:** pause | auto-forfeit | AI takeover
  - **Decision:** grace period → auto-forfeit; outcome may be game-dependent
  - **Rationale:** avoids stalled games; keeps flow simple for MVP

- **ADR-008: Host controls**
  - **Options:** full moderation suite | minimal controls | fully automatic
  - **Decision:** minimal: host can end session + replace disconnected player
  - **Rationale:** preserves hosted-session intent without building heavy moderation tools

- **ADR-009: Game UI hosting**
  - **Options:** plugin-loaded UI | separate per-game apps | bundled in main web app
  - **Decision:** bundled in main web app
  - **Rationale:** simplest deployment and integration for MVP

### Out of Scope for MVP

- Additional games beyond Connect4
- Chat
- Spectators
- Rankings, stats, or player history
- Private lobbies
- Invite links/codes
- Password reset and advanced account features

### MVP Success Criteria

- Achieves the user success metrics and KPIs defined above, especially:
  - Players can start a game within 3 minutes with ≤ 1 minute lobby waiting time
  - Connect4 matches reliably end with correct win/draw detection
  - Authoritative referee correctly rejects illegal moves

### Future Vision

- Expand game catalog beyond Connect4 (e.g., checkers) and later support more complex games (including hidden information and >2 players).
- Improve social play and rendezvous (invite links/codes, presence improvements), plus optional spectators and chat.
- Add quality-of-life features (password reset, profiles, stats/history) while keeping the platform truly free (no ads, no paid tiers).
