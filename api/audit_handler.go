package api

import (
	"net/http"
	"sort"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type AuditHandler struct {
	store *store.Store
	cfg   *config.Config
}

func NewAuditHandler(store *store.Store, cfg *config.Config) *AuditHandler {
	return &AuditHandler{store: store, cfg: cfg}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var entries interface{}
	var err error

	if userID != "" {
		entries, err = h.store.GetAuditEntriesByUser(userID)
	} else {
		entries, err = h.store.GetAllAuditEntries()
	}

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to get audit entries")
		return
	}

	if userID == "" {
		if list, ok := entries.([]*models.AuditEntry); ok {
			sort.Slice(list, func(i, j int) bool { return list[i].Timestamp < list[j].Timestamp })
			entries = list
		}
	}

	WriteJSON(w, http.StatusOK, entries)
}

func (h *AuditHandler) Verify(w http.ResponseWriter, r *http.Request) {
	result, err := auth.VerifyAuditChain(h.store)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if result.Valid {
		WriteJSON(w, http.StatusOK, map[string]interface{}{
			"chain_intact": true,
			"message":      "✅ Chain intact — no tampering detected",
			"total":        len(result.Entries),
		})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"chain_intact": false,
		"message":      "⚠️ " + result.Reason,
		"broken_at":    result.BrokenAt,
		"broken_id":    result.BrokenID,
		"reason":       result.Reason,
		"total":        len(result.Entries),
	})
}
