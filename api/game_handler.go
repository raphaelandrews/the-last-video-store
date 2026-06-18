package api

import (
	"encoding/json"
	"net/http"

	"github.com/thelastvideostore/internal/store"
)

type GameHandler struct {
	Store *store.Store
}

func (h *GameHandler) PlayStart(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	var req struct {
		GameID string `json:"game_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.GameID == "" {
		WriteError(w, http.StatusBadRequest, "game_id is required")
		return
	}

	game, err := h.Store.GetMovieByID(req.GameID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "game not found")
		return
	}

	session, err := h.Store.StartGameSession(user.ID, game.ID, game.Title, game.PlayPrice)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Play session started",
		"session": session,
		"rate":    game.PlayPrice,
	})
}

func (h *GameHandler) PlayEnd(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.SessionID == "" {
		WriteError(w, http.StatusBadRequest, "session_id is required")
		return
	}

	session, err := h.Store.EndGameSession(req.SessionID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Play session ended",
		"session": session,
	})
}

func (h *GameHandler) ActiveSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.Store.ListActiveGameSessions()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"sessions": sessions})
}

func (h *GameHandler) MySessions(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	sessions, err := h.Store.ListGameSessionsByUser(user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"sessions": sessions})
}
