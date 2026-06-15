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
	valid, err := auth.VerifyAuditChain(h.store)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"chain_intact": valid,
		"message":      chainMessage(valid),
	})
}

func chainMessage(valid bool) string {
	if valid {
		return "✅ Chain intact — no tampering detected"
	}
	return "⚠️ Chain broken — possible tampering detected"
}
