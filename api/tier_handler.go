package api

import (
	"encoding/json"
	"net/http"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type TierHandler struct {
	Store *store.Store
}

func (h *TierHandler) List(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"tiers": models.Tiers,
	})
}

func (h *TierHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req struct {
		TierName string `json:"tier_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	tier := models.TierByName(req.TierName)
	if tier == nil || tier.Name != req.TierName {
		WriteError(w, http.StatusBadRequest, "invalid tier")
		return
	}

	if err := h.Store.PurchaseTier(user.ID, tier); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "Tier purchased successfully",
		"tier":           tier.Label,
		"free_rentals":   tier.FreeRentals,
		"max_concurrent": tier.MaxConcurrent,
	})
}
