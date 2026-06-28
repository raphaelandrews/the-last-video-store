package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/store"
)

type MovieHandler struct {
	store *store.Store
	cfg   *config.Config
	hc    *crypto.HashChain
}

func NewMovieHandler(store *store.Store, cfg *config.Config, hc *crypto.HashChain) *MovieHandler {
	return &MovieHandler{store: store, cfg: cfg, hc: hc}
}

func (h *MovieHandler) List(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	mediaType := r.URL.Query().Get("media_type")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = config.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	movies, total, err := h.store.ListMoviesFiltered(genre, mediaType, offset, pageSize)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to list movies")
		return
	}

	var responses []interface{}
	for _, m := range movies {
		r := m.ToResponse()
		r.IsStaffPick = h.store.IsStaffPick(m.ID)
		responses = append(responses, r)
	}
	if responses == nil {
		responses = []interface{}{}
	}

	WriteJSON(w, http.StatusOK, MovieListResponse{
		Movies:   responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		WriteError(w, http.StatusBadRequest, "query parameter q is required")
		return
	}

	movies, err := h.store.SearchMoviesByPrefix(q, config.MaxSearchResults)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "search failed")
		return
	}

	var responses []interface{}
	for _, m := range movies {
		responses = append(responses, m.ToResponse())
	}
	if responses == nil {
		responses = []interface{}{}
	}

	WriteJSON(w, http.StatusOK, responses)
}

func (h *MovieHandler) StaffPicks(w http.ResponseWriter, r *http.Request) {
	mediaType := r.URL.Query().Get("media_type")
	ids, err := h.store.GetStaffPicks()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get staff picks")
		return
	}

	var movies []interface{}
	for _, id := range ids {
		m, err := h.store.GetMovieByID(id)
		if err != nil {
			continue
		}
		if mediaType != "" && m.MediaType != mediaType {
			continue
		}
		r := m.ToResponse()
		r.IsStaffPick = true
		movies = append(movies, r)
	}
	if movies == nil {
		movies = []interface{}{}
	}

	WriteJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) LastChance(w http.ResponseWriter, r *http.Request) {
	mediaType := r.URL.Query().Get("media_type")
	movies, err := h.store.GetLastChanceMovies()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get last chance movies")
		return
	}

	var responses []interface{}
	for _, m := range movies {
		if mediaType != "" && m.MediaType != mediaType {
			continue
		}
		responses = append(responses, m.ToResponse())
	}
	if responses == nil {
		responses = []interface{}{}
	}

	WriteJSON(w, http.StatusOK, responses)
}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movie, err := h.store.GetMovieByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "movie not found")
		return
	}

	resp := movie.ToResponse()
	resp.IsStaffPick = h.store.IsStaffPick(movie.ID)
	WriteJSON(w, http.StatusOK, resp)
}
