package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type SnackBarHandler struct {
	Store *store.Store
}

func (h *SnackBarHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSnackBarItems()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []models.SnackBarItem{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": items})
}

func (h *SnackBarHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	var req struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.ItemID == "" {
		WriteError(w, http.StatusBadRequest, "item_id is required")
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	order, err := h.Store.PlaceSnackBarOrder(user.ID, req.ItemID, req.Quantity)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Order placed successfully",
		"order":       order,
		"total_spent": order.Total,
	})
}

func (h *SnackBarHandler) Orders(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	orders, err := h.Store.ListSnackBarOrders(user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{"orders": orders})
}

func (h *SnackBarHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		Stock       int     `json:"stock"`
		Emoji       string  `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Name == "" {
		WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	item := &models.SnackBarItem{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
		Emoji:       req.Emoji,
	}

	if err := h.Store.CreateSnackBarItem(item); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]interface{}{"item": item})
}

func (h *SnackBarHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "id")

	existing, err := h.Store.GetSnackBarItem(itemID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	var req struct {
		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		Price       *float64 `json:"price,omitempty"`
		Category    *string  `json:"category,omitempty"`
		Stock       *int     `json:"stock,omitempty"`
		Emoji       *string  `json:"emoji,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Stock != nil {
		existing.Stock = *req.Stock
	}
	if req.Emoji != nil {
		existing.Emoji = *req.Emoji
	}

	if err := h.Store.UpdateSnackBarItem(existing); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{"item": existing})
}

func (h *SnackBarHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "id")

	if err := h.Store.DeleteSnackBarItem(itemID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "item deleted"})
}

func (h *SnackBarHandler) Restock(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID string `json:"item_id"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.ItemID == "" {
		WriteError(w, http.StatusBadRequest, "item_id is required")
		return
	}
	if req.Amount <= 0 {
		req.Amount = 1
	}

	item, err := h.Store.RestockSnackBarItem(req.ItemID, req.Amount)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "Item restocked successfully",
		"item":      item,
		"new_stock": item.Stock,
	})
}
