package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

type contextKey string

const (
	ctxUserKey        contextKey = "user"
	ctxPermissionsKey contextKey = "permissions"
	ctxRequestIDKey   contextKey = "request_id"
)

func newRequestID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "no-rand"
	}
	return hex.EncodeToString(b)
}

func GetUser(r *http.Request) *models.User {
	u, _ := r.Context().Value(ctxUserKey).(*models.User)
	return u
}

func GetPermissions(r *http.Request) bitmask.Permission {
	p, _ := r.Context().Value(ctxPermissionsKey).(bitmask.Permission)
	return p
}

func GetRequestID(r *http.Request) string {
	id, _ := r.Context().Value(ctxRequestIDKey).(string)
	return id
}

func AuthMiddleware(secret string, store *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				WriteError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := auth.ValidateAccessToken(tokenStr, secret)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			user, err := store.GetUserByID(claims.Subject)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "user not found")
				return
			}

			if user.Banned {
				WriteError(w, http.StatusForbidden, "account suspended")
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserKey, user)
			ctx = context.WithValue(ctx, ctxPermissionsKey, user.Tier)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(required bitmask.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perms := GetPermissions(r)
			if !bitmask.Has(perms, required) {
				WriteError(w, http.StatusForbidden, "⛔ ACCESS DENIED — Insufficient clearance")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireStaff() func(http.Handler) http.Handler {
	return RequirePermission(bitmask.PermStaff)
}

type tokenBucket struct {
	tokens   float64
	lastTime time.Time
	mu       sync.Mutex
}

func RateLimitMiddleware(rate int) func(http.Handler) http.Handler {
	buckets := &sync.Map{}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = fwd
			}

			val, _ := buckets.LoadOrStore(ip, &tokenBucket{
				tokens:   float64(rate),
				lastTime: time.Now(),
			})
			tb := val.(*tokenBucket)

			tb.mu.Lock()
			now := time.Now()
			elapsed := now.Sub(tb.lastTime).Seconds()
			tb.tokens += elapsed * float64(rate) / 60.0
			if tb.tokens > float64(rate) {
				tb.tokens = float64(rate)
			}
			tb.lastTime = now

			if tb.tokens < 1 {
				tb.mu.Unlock()
				WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			tb.tokens--
			tb.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-TOTP-Code")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = newRequestID()
			}
			w.Header().Set("X-Request-ID", id)
			ctx := context.WithValue(r.Context(), ctxRequestIDKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)

			if rec.status >= 400 {
				log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
			} else {
				log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
			}
		})
	}
}

func RecoverMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %v", err)
					WriteError(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
