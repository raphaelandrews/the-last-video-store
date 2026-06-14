package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/models"
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = config.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	movies, total, err := h.store.ListMovies(genre, offset, pageSize)
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
	movies, err := h.store.GetLastChanceMovies()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get last chance movies")
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

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Year < 1900 || req.Year > time.Now().Year()+5 {
		WriteError(w, http.StatusBadRequest, "invalid year")
		return
	}
	if req.CopiesTotal < 1 {
		req.CopiesTotal = 1
	}

	movie := &models.Movie{
		ID:              uuid.NewString(),
		Title:           req.Title,
		Year:            req.Year,
		Genre:           req.Genre,
		Format:          req.Format,
		Director:        req.Director,
		Cast:            req.Cast,
		Synopsis:        req.Synopsis,
		CopiesTotal:     req.CopiesTotal,
		CopiesAvailable: req.CopiesTotal,
		Available:       true,
		IsNewRelease:    req.IsNewRelease,
		CreatedAt:       time.Now().Unix(),
	}

	if err := h.store.CreateMovie(movie); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create movie")
		return
	}

	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionAddMovie, user.ID, movie.ID, movie.Title)

	WriteJSON(w, http.StatusCreated, movie.ToResponse())
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movie, err := h.store.GetMovieByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "movie not found")
		return
	}

	var req UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.Year != nil {
		movie.Year = *req.Year
	}
	if req.Genre != nil {
		movie.Genre = *req.Genre
	}
	if req.Format != nil {
		movie.Format = *req.Format
	}
	if req.Director != nil {
		movie.Director = *req.Director
	}
	if req.Cast != nil {
		movie.Cast = *req.Cast
	}
	if req.Synopsis != nil {
		movie.Synopsis = *req.Synopsis
	}
	if req.CopiesTotal != nil {
		movie.CopiesTotal = *req.CopiesTotal
	}
	if req.IsNewRelease != nil {
		movie.IsNewRelease = *req.IsNewRelease
	}

	if err := h.store.UpdateMovie(movie); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to update movie")
		return
	}

	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionEditMovie, user.ID, movie.ID, movie.Title)

	WriteJSON(w, http.StatusOK, movie.ToResponse())
}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movie, err := h.store.GetMovieByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "movie not found")
		return
	}

	if err := h.store.DeleteMovie(id); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to delete movie")
		return
	}

	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionDeleteMovie, user.ID, id, movie.Title)

	WriteJSON(w, http.StatusOK, SuccessResponse{Message: "movie deleted"})
}

func (h *MovieHandler) AddStaffPick(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := h.store.GetMovieByID(id); err != nil {
		WriteError(w, http.StatusNotFound, "movie not found")
		return
	}

	if err := h.store.AddStaffPick(id); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to add staff pick")
		return
	}

	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionAddStaffPick, user.ID, id, "")

	WriteJSON(w, http.StatusOK, StaffPickResponse{StaffPick: true})
}

func (h *MovieHandler) RemoveStaffPick(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.store.RemoveStaffPick(id); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to remove staff pick")
		return
	}

	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionRemoveStaffPick, user.ID, id, "")

	WriteJSON(w, http.StatusOK, StaffPickResponse{StaffPick: false})
}
