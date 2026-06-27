package pages

import (
	"github.com/thelastvideostore/internal/models"
)

type NavigateMsg struct{ Page string }
type ErrorMsg struct{ Message string }

type LoginRequestMsg struct {
	Username string
	Password string
}

type LoginSuccessMsg struct {
	AccessToken  string
	RefreshToken string
	User         *models.UserResponse
}

type RegisterRequestMsg struct {
	Username string
	Password string
}
