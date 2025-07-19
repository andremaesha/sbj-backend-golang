package usecase

import (
	"context"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/helpers"
	"time"
)

type authUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) web.AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		contextTimeout: contextTimeout,
	}
}

// GetUserFromSession retrieves a user from a session ID
func (au *authUsecase) GetUserFromSession(ctx context.Context, sessionID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, au.contextTimeout)
	defer cancel()

	return au.userRepository.GetSession(ctx, sessionID)
}

// DecryptSessionID decrypts a session cookie to get the session ID
func (au *authUsecase) DecryptSessionID(key, sessionCookie string) (string, error) {
	return helpers.DecryptAES(sessionCookie, key)
}
