package api

import (
	"net/http"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type InventoryHandler struct {
	Store *store.Store
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	items, err := h.Store.ListInventory(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []models.InventoryItem{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}
