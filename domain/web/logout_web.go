package web

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type LogoutResponse struct {
	Message string `json:"message"`
}

type LogoutUsecase interface {
	DeleteSession(c context.Context, idSession string) error
	DecryptSession(key, data string) string
	ValidateSession(sessionId string) error
	CreateExpiredCookie() *fiber.Cookie
}
