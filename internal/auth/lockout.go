package auth

import (
	"fmt"
	"time"

	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/store"
)

var ErrAccountLocked = fmt.Errorf("account is locked due to too many failed attempts")

func CheckLoginAttempts(s *store.Store, username string) error {
	locked, err := s.IsUserLocked(username)
	if err != nil {
		return fmt.Errorf("check login attempts: %w", err)
	}
	if locked {
		return ErrAccountLocked
	}
	return nil
}

func RecordFailedAttempt(s *store.Store, username string) error {
	count, err := s.IncrementFailedAttempts(username)
	if err != nil {
		return fmt.Errorf("record failed attempt: %w", err)
	}

	if count >= config.MaxLoginAttempts {
		until := time.Now().Add(config.LockoutDuration)
		if err := s.LockUserUntil(username, until.Unix()); err != nil {
			return fmt.Errorf("lock user: %w", err)
		}
	}

	return nil
}

func RecordSuccessfulLogin(s *store.Store, username string) error {
	if err := s.ResetFailedAttempts(username); err != nil {
		return fmt.Errorf("record success: %w", err)
	}
	return nil
}
