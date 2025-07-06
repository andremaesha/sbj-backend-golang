package web

import (
	"context"
	"sbj-backend/domain"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type LoginUsecase interface {
	GetUserByEmail(c context.Context, email string) (*domain.User, error)
	ValidateUserVerified(verified bool) error
	SetSession(c context.Context, expire int, idSession string, data *domain.User) error
	EncryptSession(key, data string) string
}
