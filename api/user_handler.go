package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type UserHandler struct {
	store *store.Store
	cfg   *config.Config
	hc    *crypto.HashChain
}

func NewUserHandler(store *store.Store, cfg *config.Config, hc *crypto.HashChain) *UserHandler {
	return &UserHandler{store: store, cfg: cfg, hc: hc}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	var responses []interface{}
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	WriteJSON(w, http.StatusOK, responses)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	admin := GetUser(r)

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if h.store.UserExists(req.Username) {
		WriteError(w, http.StatusConflict, "username already taken")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	now := time.Now().Unix()
	tier := parseTier(req.Tier)
	if tier == 0 {
		tier = bitmask.TierBronze
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		PasswordHash: hash,
		Tier:         tier,
		MaxRentals:   bitmask.MaxRentalsForTier(tier),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.store.CreateUser(user); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionRegister, admin.ID, user.ID, user.Username)

	WriteJSON(w, http.StatusCreated, user.ToResponse())
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	admin := GetUser(r)

	target, err := h.store.GetUserByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Tier != nil {
		oldTier := bitmask.TierName(target.Tier)
		newTier := parseTier(*req.Tier)
		if newTier == 0 {
			WriteError(w, http.StatusBadRequest, "invalid tier")
			return
		}

		if !admin.CanAdmin() && newTier > admin.Tier {
			WriteError(w, http.StatusForbidden, "cannot promote beyond your own tier")
			return
		}

		target.Tier = newTier
		target.MaxRentals = bitmask.MaxRentalsForTier(newTier)

		action := models.ActionPromote
		if newTier < target.Tier {
			action = models.ActionDemote
		}
		auth.AppendAuditEntry(h.store, h.hc, action, admin.ID, target.ID,
			oldTier+" → "+bitmask.TierName(newTier))
	}

	if req.Banned != nil {
		target.Banned = *req.Banned
		action := models.ActionBan
		if !*req.Banned {
			action = models.ActionUnban
		}
		auth.AppendAuditEntry(h.store, h.hc, action, admin.ID, target.ID, "")
	}

	target.UpdatedAt = time.Now().Unix()

	if err := h.store.UpdateUser(target); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	WriteJSON(w, http.StatusOK, target.ToResponse())
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	admin := GetUser(r)

	if err := h.store.DeleteUser(id); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, "delete_user", admin.ID, id, "")

	WriteJSON(w, http.StatusOK, SuccessResponse{Message: "user deleted"})
}

func (h *UserHandler) TOTPSetup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	currentUser := GetUser(r)

	if id != currentUser.ID && !currentUser.CanAdmin() {
		WriteError(w, http.StatusForbidden, "⛔ ACCESS DENIED")
		return
	}

	user, err := h.store.GetUserByID(id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	var req TOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Enabled {
		secret, err := auth.GenerateTOTPSecret()
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to generate secret")
			return
		}

		encrypted, err := crypto.Encrypt([]byte(secret), []byte(h.cfg.AESKey))
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to encrypt secret")
			return
		}

		if err := h.store.SaveTOTPSecret(user.ID, encrypted); err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to save secret")
			return
		}

		user.TOTPEnabled = true
		h.store.UpdateUser(user)

		auth.AppendAuditEntry(h.store, h.hc, models.ActionTOTPEnabled, currentUser.ID, user.ID, "")

		url := auth.GenerateTOTPURL("TLVS", user.Username, secret)
		WriteJSON(w, http.StatusOK, TOTPSetupResponse{Secret: secret, URL: url})
	} else {
		h.store.DeleteTOTPSecret(user.ID)
		user.TOTPEnabled = false
		h.store.UpdateUser(user)

		auth.AppendAuditEntry(h.store, h.hc, models.ActionTOTPDisabled, currentUser.ID, user.ID, "")

		WriteJSON(w, http.StatusOK, SuccessResponse{Message: "totp disabled"})
	}
}

func parseTier(s string) bitmask.Permission {
	switch s {
	case "bronze":
		return bitmask.TierBronze
	case "silver":
		return bitmask.TierSilver
	case "gold":
		return bitmask.TierGold
	case "employee":
		return bitmask.TierEmployee
	case "supervisor":
		return bitmask.TierSupervisor
	case "manager":
		return bitmask.TierManager
	case "owner":
		return bitmask.TierOwner
	default:
		return 0
	}
}
