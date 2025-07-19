package web

import (
	"context"
	"sbj-backend/domain"
)

// AuthUsecase defines the interface for authentication operations
type AuthUsecase interface {
	// GetUserFromSession retrieves a user from a session ID
	GetUserFromSession(ctx context.Context, sessionID string) (*domain.User, error)

	// DecryptSessionID decrypts a session cookie to get the session ID
	DecryptSessionID(key, sessionCookie string) (string, error)
}
