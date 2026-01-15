package lobby_service

import (
	"net/http"

	"mygameplatform/internal/httpx"
)

type gamesResponse struct {
	Games []gameInfo `json:"games"`
}

type gameInfo struct {
	GameID       string `json:"game_id"`
	Name         string `json:"name"`
	MinPlayers   int    `json:"min_players"`
	MaxPlayers   int    `json:"max_players"`
	RulesSummary string `json:"rules_summary"`
}

func (h *Handler) handleGames(w http.ResponseWriter, r *http.Request) {
	resp := gamesResponse{
		Games: []gameInfo{
			{
				GameID:       "connect4",
				Name:         "Connect4",
				MinPlayers:   2,
				MaxPlayers:   2,
				RulesSummary: "7x6 board; gravity drop; 4-in-a-row wins",
			},
			{
				GameID:       "draughts_10x10",
				Name:         "Draughts (10x10)",
				MinPlayers:   2,
				MaxPlayers:   2,
				RulesSummary: "International Draughts: forced captures, max-capture priority, multi-jumps, flying kings",
			},
		},
	}

	h.logRequest(r, http.StatusOK, "", nil)
	httpx.WriteJSON(w, http.StatusOK, resp)
}
