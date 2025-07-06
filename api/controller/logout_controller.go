package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain/web"
	"time"
)

type LogoutController struct {
	LogoutUsecase web.LogoutUsecase
	Env           *bootstrap.Env
	Session       *session.Store
}

func (lc *LogoutController) Logout(c *fiber.Ctx) error {
	sessionId := c.Cookies("session_id")
	println(sessionId)
	if sessionId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&web.LogoutResponse{
			Message: "session not found",
		})
	}

	content := lc.LogoutUsecase.DecryptSession(lc.Env.Key, sessionId)

	err := lc.LogoutUsecase.DeleteSession(c.Context(), content)
	if err != nil {
		panic(err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(web.LogoutResponse{
		Message: "Logout successful",
	})
}
