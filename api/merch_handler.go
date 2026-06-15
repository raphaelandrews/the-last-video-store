package api

import (
	"encoding/json"
	"net/http"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type MerchHandler struct {
	Store *store.Store
}

func (h *MerchHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMerchItems()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []models.MerchItem{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *MerchHandler) Redeem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		ItemID string `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.ItemID == "" {
		WriteError(w, http.StatusBadRequest, "item_id is required")
		return
	}

	item, err := h.Store.RedeemMerchItem(req.ItemID, userID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":      "Redemption successful",
		"item":         item,
		"points_spent": item.PointsCost,
	})
}
