package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type AuthHandler struct {
	store *store.Store
	cfg   *config.Config
	hc    *crypto.HashChain
}

func NewAuthHandler(store *store.Store, cfg *config.Config, hc *crypto.HashChain) *AuthHandler {
	return &AuthHandler{store: store, cfg: cfg, hc: hc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 || len(req.Username) > 20 {
		WriteError(w, http.StatusBadRequest, "username must be 3–20 characters")
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username) {
		WriteError(w, http.StatusBadRequest, "username must be alphanumeric")
		return
	}
	if len(req.Password) < 6 {
		WriteError(w, http.StatusBadRequest, "password must be at least 6 characters")
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
	tier := bitmask.TierBronze
	maxRentals := bitmask.MaxRentalsForTier(tier)

	user := &models.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		PasswordHash: hash,
		Tier:         tier,
		MaxRentals:   maxRentals,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.store.CreateUser(user); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionRegister, user.ID, "", "new user registered")
	WriteJSON(w, http.StatusCreated, user.ToResponse())
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := auth.CheckLoginAttempts(h.store, req.Username); err != nil {
		WriteError(w, http.StatusTooManyRequests, "account locked — too many failed attempts")
		return
	}

	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		auth.RecordFailedAttempt(h.store, req.Username)
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		auth.RecordFailedAttempt(h.store, req.Username)
		WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	auth.RecordSuccessfulLogin(h.store, req.Username)

	if user.Banned {
		WriteError(w, http.StatusForbidden, "account suspended")
		return
	}

	if user.TOTPEnabled {
		tempToken, expiresAt, err := auth.GenerateTOTPTempToken(user.ID, h.cfg.JWTSecret)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to generate token")
			return
		}
		WriteJSON(w, http.StatusOK, map[string]interface{}{
			"totp_required": true,
			"temp_token":    tempToken,
			"expires_at":    expiresAt,
		})
		return
	}

	h.issueTokens(w, r, user)
}

func (h *AuthHandler) LoginTOTP(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		WriteError(w, http.StatusUnauthorized, "missing temp token")
		return
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	userID, err := auth.ValidateTOTPTempToken(tokenStr, h.cfg.JWTSecret)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "invalid or expired temp token")
		return
	}

	var req TOTPLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "user not found")
		return
	}

	encryptedSecret, err := h.store.GetTOTPSecret(user.ID)
	if err != nil || len(encryptedSecret) == 0 {
		WriteError(w, http.StatusInternalServerError, "totp not configured")
		return
	}

	decrypted, err := crypto.Decrypt(encryptedSecret, []byte(h.cfg.AESKey))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to decrypt totp secret")
		return
	}

	if !auth.ValidateTOTPCode(string(decrypted), req.Code) {
		count, _ := h.store.IncrementTOTPFailures(user.ID)
		if count >= 3 {
			h.store.LockTOTPUserUntil(user.ID, time.Now().Add(10*time.Minute).Unix())
		}
		WriteError(w, http.StatusUnauthorized, "invalid totp code")
		return
	}

	h.store.ResetTOTPFailures(user.ID)
	h.issueTokens(w, r, user)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken, h.cfg.JWTSecret)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	valid, _ := h.store.ValidateRefreshToken(claims.Subject, claims.TokenID)
	if !valid {
		WriteError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}

	h.store.InvalidateRefreshToken(claims.Subject, claims.TokenID)

	user, err := h.store.GetUserByID(claims.Subject)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "user not found")
		return
	}

	pair, err := auth.GenerateTokenPair(user.ID, uint16(user.Tier), h.cfg.JWTSecret)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	refreshClaims, _ := auth.ValidateRefreshToken(pair.RefreshToken, h.cfg.JWTSecret)
	if refreshClaims != nil {
		h.store.SaveRefreshToken(user.ID, refreshClaims.TokenID, refreshClaims.ExpiresAt.Unix())
	}

	WriteJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.ExpiresAt,
		User:         user.ToResponse(),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	auth.AppendAuditEntry(h.store, h.hc, models.ActionLogout, user.ID, "", "logout")
	WriteJSON(w, http.StatusOK, SuccessResponse{Message: "logged out"})
}

func (h *AuthHandler) issueTokens(w http.ResponseWriter, r *http.Request, user *models.User) {
	pair, err := auth.GenerateTokenPair(user.ID, uint16(user.Tier), h.cfg.JWTSecret)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to generate tokens")
		return
	}

	refreshClaims, _ := auth.ValidateRefreshToken(pair.RefreshToken, h.cfg.JWTSecret)
	if refreshClaims != nil {
		h.store.SaveRefreshToken(user.ID, refreshClaims.TokenID, refreshClaims.ExpiresAt.Unix())
	}

	auth.AppendAuditEntry(h.store, h.hc, models.ActionLogin, user.ID, "", "login")

	WriteJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.ExpiresAt,
		User:         user.ToResponse(),
	})
}
