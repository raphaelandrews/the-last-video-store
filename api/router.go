package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/store"
)

func NewRouter(store *store.Store, cfg *config.Config, hc *crypto.HashChain) http.Handler {
	r := chi.NewRouter()

	r.Use(RecoverMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(LoggingMiddleware())
	r.Use(CORSMiddleware())
	r.Use(RateLimitMiddleware(config.RateLimitPerMinute))

	authHandler := NewAuthHandler(store, cfg, hc)
	movieHandler := NewMovieHandler(store, cfg, hc)
	rentalHandler := NewRentalHandler(store, cfg, hc)
	userHandler := NewUserHandler(store, cfg, hc)
	wishlistHandler := NewWishlistHandler(store, cfg, hc)
	auditHandler := NewAuditHandler(store, cfg)
	merchHandler := &MerchHandler{Store: store}
	inventoryHandler := &InventoryHandler{Store: store}
	tierHandler := &TierHandler{Store: store}
	snackBarHandler := &SnackBarHandler{Store: store}
	gameHandler := &GameHandler{Store: store}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/login/totp", authHandler.LoginTOTP)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(cfg.JWTSecret, store))

			r.Post("/auth/refresh", authHandler.Refresh)
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/me", authHandler.Me)

			r.Route("/movies", func(r chi.Router) {
				r.Get("/", movieHandler.List)
				r.Get("/search", movieHandler.Search)
				r.Get("/staff-picks", movieHandler.StaffPicks)
				r.Get("/last-chance", movieHandler.LastChance)
				r.Get("/{id}", movieHandler.GetByID)

				r.Group(func(r chi.Router) {
					r.Use(RequirePermission(bitmask.PermAdmin))
					r.Post("/", movieHandler.Create)
					r.Put("/{id}", movieHandler.Update)
					r.Delete("/{id}", movieHandler.Delete)
					r.Post("/{id}/staff-pick", movieHandler.AddStaffPick)
					r.Delete("/{id}/staff-pick", movieHandler.RemoveStaffPick)
				})
			})

			r.Route("/rentals", func(r chi.Router) {
				r.Post("/rent", rentalHandler.Rent)
				r.Post("/return", rentalHandler.Return)
				r.Post("/extend", rentalHandler.Extend)
				r.Get("/history", rentalHandler.History)
			})

			r.Route("/wishlist", func(r chi.Router) {
				r.Get("/", wishlistHandler.List)
				r.Post("/", wishlistHandler.Add)
				r.Delete("/{movieID}", wishlistHandler.Remove)
				r.Get("/check/{movieID}", wishlistHandler.Check)
			})

			r.Route("/users", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(RequirePermission(bitmask.PermManageUsers))
					r.Get("/", userHandler.List)
					r.Post("/", userHandler.Create)
					r.Put("/{id}", userHandler.Update)
				})

				r.Group(func(r chi.Router) {
					r.Use(RequirePermission(bitmask.PermAdmin))
					r.Delete("/{id}", userHandler.Delete)
				})

				r.Post("/{id}/totp", userHandler.TOTPSetup)
				r.Post("/me/topup", userHandler.TopUp)
			})

			r.Route("/audit", func(r chi.Router) {
				r.Use(RequirePermission(bitmask.PermManageUsers))
				r.Get("/", auditHandler.List)
				r.Get("/verify", auditHandler.Verify)
			})

			r.Route("/merch", func(r chi.Router) {
				r.Get("/", merchHandler.List)
				r.Post("/redeem", merchHandler.Redeem)
			})

			r.Route("/inventory", func(r chi.Router) {
				r.Get("/", inventoryHandler.List)
			})

			r.Route("/tiers", func(r chi.Router) {
				r.Get("/", tierHandler.List)
				r.Post("/purchase", tierHandler.Purchase)
			})

			r.Route("/snackbar", func(r chi.Router) {
				r.Get("/", snackBarHandler.List)
				r.Post("/order", snackBarHandler.PlaceOrder)
				r.Get("/orders", snackBarHandler.Orders)

				r.Group(func(r chi.Router) {
					r.Use(RequirePermission(bitmask.PermSnackBarManage))
					r.Post("/items", snackBarHandler.CreateItem)
					r.Put("/items/{id}", snackBarHandler.UpdateItem)
					r.Delete("/items/{id}", snackBarHandler.DeleteItem)
					r.Post("/restock", snackBarHandler.Restock)
				})
			})

			r.Route("/games", func(r chi.Router) {
				r.Get("/my-sessions", gameHandler.MySessions)
				r.Post("/play/start", gameHandler.PlayStart)
				r.Post("/play/end", gameHandler.PlayEnd)

				r.Group(func(r chi.Router) {
					r.Use(RequirePermission(bitmask.PermGameManage))
					r.Get("/play/active", gameHandler.ActiveSessions)
				})
			})
		})
	})

	return r
}
