package api

import (
	"net/http"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
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

	WriteJSON(w, http.StatusOK, entries)
}

func (h *AuditHandler) Verify(w http.ResponseWriter, r *http.Request) {
	valid, chainErr := auth.VerifyAuditChain(h.store)
	if chainErr != nil && !valid {
		// chainErr can only be non-nil when valid is false here, but we
		// guard with the check for clarity.
		WriteJSON(w, http.StatusOK, map[string]interface{}{
			"chain_intact": false,
			"message":      "⚠️ " + chainErr.Error(),
			"broken_at":    chainErr.BrokenAt,
			"reason":       chainErr.Reason,
		})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"chain_intact": true,
		"message":      "✅ Chain intact — no tampering detected",
	})
}
