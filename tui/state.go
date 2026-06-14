package tui

import (
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
)

type SessionState struct {
	AccessToken  string
	RefreshToken string
	User         *models.UserResponse
	Permissions  bitmask.Permission
	IsLoggedIn   bool
	APIBaseURL   string
	cache        map[string]interface{}
}

func NewSessionState(apiURL string) *SessionState {
	return &SessionState{
		APIBaseURL: apiURL,
		cache:      make(map[string]interface{}),
	}
}

func (s *SessionState) Login(accessToken, refreshToken string, user *models.UserResponse) {
	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	s.User = user
	s.Permissions = user.Tier
	s.IsLoggedIn = true
}

func (s *SessionState) Logout() {
	s.AccessToken = ""
	s.RefreshToken = ""
	s.User = nil
	s.Permissions = 0
	s.IsLoggedIn = false
	s.cache = make(map[string]interface{})
}

func (s *SessionState) HasPermission(perm bitmask.Permission) bool {
	return bitmask.Has(s.Permissions, perm)
}

func (s *SessionState) CanAccessAdmin() bool {
	return s.HasPermission(bitmask.PermManageUsers)
}

func (s *SessionState) CanAdminMovies() bool {
	return s.HasPermission(bitmask.PermAdmin)
}

func (s *SessionState) IsStaff() bool {
	return s.HasPermission(bitmask.PermStaff)
}

func (s *SessionState) CacheGet(key string) (interface{}, bool) {
	v, ok := s.cache[key]
	return v, ok
}

func (s *SessionState) CacheSet(key string, val interface{}) {
	s.cache[key] = val
}

func (s *SessionState) CacheClear() {
	s.cache = make(map[string]interface{})
}
