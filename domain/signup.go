package domain

import "context"

type SignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Avatar    string `json:"avatar"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

type SignupUsecase interface {
	Create(c context.Context, user *User) error
	GetUserByEmail(c context.Context, email string) (*User, error)
}
