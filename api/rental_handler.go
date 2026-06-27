package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type RentalHandler struct {
	store *store.Store
	cfg   *config.Config
	hc    *crypto.HashChain
}

func NewRentalHandler(store *store.Store, cfg *config.Config, hc *crypto.HashChain) *RentalHandler {
	return &RentalHandler{store: store, cfg: cfg, hc: hc}
}

func (h *RentalHandler) Rent(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req RentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	movie, err := h.store.GetMovieByID(req.MovieID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "movie not found")
		return
	}

	if !movie.HasCopies() {
		WriteError(w, http.StatusConflict, "no copies available — join waitlist?")
		return
	}

	if movie.IsNewRelease && !user.CanReserve() {
		WriteError(w, http.StatusForbidden, "gold plan required for new releases")
		return
	}

	freeRental := false
	if req.UseTicket && user.FreeRentals > 0 {
		user.FreeRentals--
		freeRental = true
	} else if user.AtRentalLimit() {
		WriteError(w, http.StatusForbidden, "rental limit reached — use a ticket, return a movie, or upgrade your tier")
		return
	}

	if !freeRental {
		cost := models.MovieCost(movie.RentalPrice, movie.Format)
		if user.Balance < cost {
			WriteError(w, http.StatusPaymentRequired, fmt.Sprintf("insufficient balance: need $%.2f, have $%.2f", cost, user.Balance))
			return
		}
		user.Balance -= cost
		user.RentalCount++
	}

	now := time.Now().Unix()
	rental := &models.Rental{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		MovieID:      movie.ID,
		MovieFormat:  movie.Format,
		RentedAt:     now,
		DueDate:      models.DueDateForFormat(movie.Format, now),
		Status:       models.RentalActive,
		IsFreeRental: freeRental,
	}

	if movie.Format == models.FormatVHS {
		rental.NeedsRewind = rand.Intn(100) < 30
	}

	movie.CopiesAvailable--
	if movie.CopiesAvailable == 0 {
		movie.Available = false
	}

	if err := h.store.UpdateMovie(movie); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to update movie")
		return
	}

	if err := h.store.UpdateUser(user); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	if err := h.store.CreateRental(rental); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create rental")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionRent, user.ID, movie.ID, movie.Title)

	resp := rental.ToResponse(movie.Title)
	WriteJSON(w, http.StatusCreated, resp)
}

func (h *RentalHandler) Return(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req ReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rental, err := h.store.GetRentalByID(req.RentalID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "rental not found")
		return
	}

	if rental.Status == models.RentalReturned {
		WriteError(w, http.StatusConflict, "rental already returned")
		return
	}

	isOwner := rental.UserID == user.ID
	isStaff := user.HasStaffAccess()
	if !isOwner && !isStaff {
		WriteError(w, http.StatusForbidden, "⛔ ACCESS DENIED — not your rental")
		return
	}

	now := time.Now().Unix()
	rental.ReturnedAt = now
	rental.Status = models.RentalReturned

	rental.LateFee = rental.CalculateLateFee(now)
	rental.RewindFee = rental.CalculateRewindFee()

	movie, err := h.store.GetMovieByID(rental.MovieID)
	if err == nil {
		movie.CopiesAvailable++
		if movie.CopiesAvailable > movie.CopiesTotal {
			movie.CopiesAvailable = movie.CopiesTotal
		}
		movie.Available = movie.CopiesAvailable > 0
		h.store.UpdateMovie(movie)
	}

	rentalUser, err := h.store.GetUserByID(rental.UserID)
	if err == nil {
		rentalUser.RentalCount--
		if rentalUser.RentalCount < 0 {
			rentalUser.RentalCount = 0
		}
		pointsEarned := 0
		if rental.LateFee == 0 && rental.RewindFee == 0 {
			rentalUser.PopcornPoints += 10
			pointsEarned += 10
		} else if rental.LateFee > 0 {
			rentalUser.PopcornPoints -= 5
			pointsEarned -= 5
		}
		rentalUser.Balance -= rental.LateFee + rental.RewindFee
		inventory, _ := h.store.ListInventory(rentalUser.ID)
		for _, item := range inventory {
			if item.MerchID == "merch-popcorn-bucket" {
				rentalUser.PopcornPoints += 5
				pointsEarned += 5
				break
			}
		}
		rental.PointsEarned = pointsEarned
		h.store.UpdateUser(rentalUser)
		h.store.UpdateRental(rental)
		auth.AppendAuditEntry(h.store, h.hc, models.ActionReturn, user.ID, rental.MovieID,
			h.movieTitle(rental.MovieID))
		resp := rental.ToResponse(h.movieTitle(rental.MovieID))
		WriteJSON(w, http.StatusOK, resp)
		return
	}

	h.store.UpdateRental(rental)

	auth.AppendAuditEntry(h.store, h.hc, models.ActionReturn, user.ID, rental.MovieID,
		h.movieTitle(rental.MovieID))

	resp := rental.ToResponse(h.movieTitle(rental.MovieID))
	WriteJSON(w, http.StatusOK, resp)
}

func (h *RentalHandler) History(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	rentals, err := h.store.GetRentalHistoryByUser(user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get rental history")
		return
	}

	var responses []models.RentalResponse
	for _, rental := range rentals {
		resp := rental.ToResponse(h.movieTitle(rental.MovieID))
		responses = append(responses, resp)
	}
	if responses == nil {
		responses = []models.RentalResponse{}
	}

	WriteJSON(w, http.StatusOK, responses)
}

func (h *RentalHandler) Extend(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req struct {
		RentalID string `json:"rental_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.RentalID == "" {
		WriteError(w, http.StatusBadRequest, "rental_id is required")
		return
	}

	const extendMinutes = 1
	const cost = 30

	if err := h.store.ExtendRental(req.RentalID, user.ID, extendMinutes, cost); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":  fmt.Sprintf("Extended by %d min for %d 🍿", extendMinutes, cost),
		"extended": extendMinutes,
	})
}

func (h *RentalHandler) movieTitle(id string) string {
	m, err := h.store.GetMovieByID(id)
	if err != nil {
		return id
	}
	return m.Title
}
