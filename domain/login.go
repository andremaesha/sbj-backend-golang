package domain

import "context"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type LoginUsecase interface {
	GetUserByEmail(c context.Context, email string) (*User, error)
}
