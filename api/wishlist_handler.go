package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type WishlistHandler struct {
	store *store.Store
	cfg   *config.Config
	hc    *crypto.HashChain
}

func NewWishlistHandler(store *store.Store, cfg *config.Config, hc *crypto.HashChain) *WishlistHandler {
	return &WishlistHandler{store: store, cfg: cfg, hc: hc}
}

func (h *WishlistHandler) List(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	items, err := h.store.GetWishlist(user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get wishlist")
		return
	}

	type WishlistEntry struct {
		ID        string `json:"id"`
		MovieID   string `json:"movie_id"`
		Title     string `json:"title"`
		Available bool   `json:"available"`
		AddedAt   int64  `json:"added_at"`
	}

	var entries []WishlistEntry
	for _, item := range items {
		entry := WishlistEntry{
			ID:      item.ID,
			MovieID: item.MovieID,
			AddedAt: item.AddedAt,
		}
		movie, err := h.store.GetMovieByID(item.MovieID)
		if err == nil {
			entry.Title = movie.Title
			entry.Available = movie.HasCopies()
		}
		entries = append(entries, entry)
	}
	if entries == nil {
		entries = []WishlistEntry{}
	}

	WriteJSON(w, http.StatusOK, entries)
}

func (h *WishlistHandler) Add(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req struct {
		MovieID string `json:"movie_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inWishlist, _ := h.store.IsInWishlist(user.ID, req.MovieID)
	if inWishlist {
		WriteError(w, http.StatusConflict, "already in wishlist")
		return
	}

	if err := h.store.AddToWishlist(user.ID, req.MovieID); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to add to wishlist")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionAddToWishlist, user.ID, req.MovieID, "")

	WriteJSON(w, http.StatusCreated, SuccessResponse{Message: "added to wishlist"})
}

func (h *WishlistHandler) Remove(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	movieID := chi.URLParam(r, "movieID")

	if err := h.store.RemoveFromWishlist(user.ID, movieID); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to remove from wishlist")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionRemoveFromWishlist, user.ID, movieID, "")

	WriteJSON(w, http.StatusOK, SuccessResponse{Message: "removed from wishlist"})
}

func (h *WishlistHandler) Check(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	movieID := chi.URLParam(r, "movieID")

	inWishlist, _ := h.store.IsInWishlist(user.ID, movieID)

	WriteJSON(w, http.StatusOK, map[string]bool{"in_wishlist": inWishlist})
}
