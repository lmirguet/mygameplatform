# Story 2.7: Responsive UI for games + lobby list + join/create + copy invite

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a signed-in user,
I want a responsive UI to select a game and join/create a lobby (and copy an invite link),
so that I can quickly get into a lobby on desktop or mobile.

## Acceptance Criteria

1. **Given** I am authenticated **When** I open the Games screen **Then** I see Connect4 as available **And** I see Draughts (10x10) as available.
2. **Given** I select a game **When** I view the lobby list **Then** I see the list of public lobbies for that game with seat availability **And** I can join an open lobby or create a new lobby.
3. **Given** I am in a lobby **When** I click “Copy invite link” **Then** the lobby URL is copied **And** I get a clear confirmation message.
4. **Given** I use a mobile browser **When** I view the lobby list UI **Then** it remains usable and readable (responsive layout).

## Tasks / Subtasks

- [ ] Build Games + Lobby List screens in web app (AC: 1-4)
  - [ ] Add Games view that lists Connect4 and Draughts (10x10)
  - [ ] Add Lobby List view for selected game with seat availability
  - [ ] Add Join and Create actions (wire to APIs from Stories 2.3–2.5)
  - [ ] Add Copy Invite Link button with confirmation toast/banner
  - [ ] Ensure responsive layout for mobile and desktop
- [ ] Integrate with backend endpoints (AC: 1-3)
  - [ ] `GET /api/v1/games`
  - [ ] `GET /api/v1/lobbies?game={game_id}`
  - [ ] `POST /api/v1/lobbies/{lobby_id}/join`
  - [ ] `POST /api/v1/lobbies` (create)
  - [ ] `POST /api/v1/lobbies/auto-join`
  - [ ] `GET /api/v1/lobbies/{lobby_id}` (for share URL if needed)
- [ ] Add UI tests or component tests where applicable (AC: 1-4)
  - [ ] Games list renders both games
  - [ ] Lobby list renders seat availability and join/create actions
  - [ ] Copy invite link triggers clipboard call and confirmation state

## Dev Notes

- Source of truth: `/_bmad-output/planning-artifacts/epics.md` (Story 2.7).
- UX spec calls for a “light table” lobby list layout with responsive stacking on mobile.
- Use Tailwind tokens and responsive utilities from UX spec.
- Keep copy invite link flow simple; no modal required.

### Project Structure Notes

- Web app routes: `web/src/app/routes.tsx`.
- Games/Lobbies pages under `web/src/pages/`.
- Shared UI components under `web/src/components/`.
- API clients under `web/src/api/`.

### References

- Story 2.7 definition and ACs: `/_bmad-output/planning-artifacts/epics.md#Story-2-7-Responsive-UI-for-games--lobby-list--joincreate--copy-invite`
- UX design system and layout: `/_bmad-output/planning-artifacts/ux-design-specification.md#Design-Direction-Decision`
- API conventions: `/_bmad-output/planning-artifacts/architecture.md#Implementation-Patterns--Consistency-Rules`

## Dev Agent Record

### Agent Model Used

GPT-5 (Codex CLI)

### Debug Log References

### Completion Notes List

- Ultimate context engine analysis completed - comprehensive developer guide created

### File List

- `/_bmad-output/implementation-artifacts/2-7-responsive-ui-for-games-lobby-list-joincreate-copy-invite.md`

## Developer Context

### Story Foundation (from Epic 2)

- Epic: Game Discovery, Lobby Join/Create, and Invite Links.
- Story 2.7: Responsive UI for games + lobby list + join/create + copy invite.
- FRs: FR4, FR5, FR6, FR7, FR8, FR14, FR18, FR19.

### Business Context

- Lobby UI is the core “speed to seat” experience and must be simple and fast.
- Visual clarity and seat availability directly impact conversion.

### UX Context

- Use the “light table” lobby list layout defined in UX spec.
- Mobile-first; lobby table stacks rows on small screens.
- Provide clear feedback on actions (join/create/copy invite).

## Technical Requirements

- Frontend routes:
  - `/games` (games list)
  - `/lobbies?game={game_id}` (lobby list by game)
  - `/lobbies/:lobby_id` (lobby detail/invite context)
- Games list:
  - Fetch from `GET /api/v1/games`.
  - Show Connect4 and Draughts (10x10).
- Lobby list:
  - Fetch from `GET /api/v1/lobbies?game={game_id}`.
  - Show `game_id`, `max_seats`, `occupied_seats`, `seats_available`.
  - Join button enabled only when `seats_available > 0`.
  - Create button available when no open lobbies (or as secondary action).
- Actions:
  - Join: `POST /api/v1/lobbies/{lobby_id}/join`.
  - Create: `POST /api/v1/lobbies`.
  - Auto-join (optional UX shortcut): `POST /api/v1/lobbies/auto-join`.
  - Copy invite: use `share_url` from lobby detail or list; call clipboard API and show confirmation.
- Error handling:
  - Show friendly inline error for “lobby full” and auth failures.
  - If unauthenticated, redirect to login with return path preserved.
- Styling:
  - Tailwind utility classes; use tokens defined in UX spec (primary blue, teal accents, neutral background).
  - Ensure readable on mobile and desktop; 8px spacing grid.

## Architecture Compliance

- JSON fields are `snake_case`.
- Use REST endpoints from `/api/v1` only.
- Do not hardcode hostnames in `share_url` usage.

## Library / Framework Requirements

- React + Vite + TypeScript (existing web app).
- Tailwind CSS for styling.
- React Router for routing.

## File Structure Requirements

- Routes: `web/src/app/routes.tsx`.
- Pages: `web/src/pages/` (games, lobbies, lobby detail).
- API client: `web/src/api/`.
- Shared UI components: `web/src/components/`.

## Testing Requirements

- Games list renders both games.
- Lobby list renders seat availability and correct button states.
- Join/Create actions call correct endpoints.
- Copy invite uses clipboard and shows confirmation.
- Mobile layout stacks lobby rows.

## Previous Story Intelligence

- Stories 2.1–2.6 define backend endpoints; UI must align with those contracts.
- Lobby `share_url` is a stable path `/lobbies/{lobby_id}`.
- Keep UI errors aligned with backend error codes.

## Git Intelligence Summary

- Existing web app is Vite + React + TS; App styles updated during Epic 1.
- No lobby UI yet; implement new pages/components under `web/src/`.

## Latest Tech Information

- Tailwind CSS and React Router versions are defined in `web/package.json`; follow the lockfile.

## Project Context Reference

- No `project-context.md` found in repo at time of story creation.

## Story Completion Status

- Status set to **ready-for-dev**.
- Completion note added: “Ultimate context engine analysis completed - comprehensive developer guide created”.
