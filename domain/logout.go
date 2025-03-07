package domain

import "context"

type LogoutResponse struct {
	Message string `json:"message"`
}

type LogoutUsecase interface {
	DeleteSession(c context.Context, idSession string) error
	DecryptSession(key, data string) string
}
