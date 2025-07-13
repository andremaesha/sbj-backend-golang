package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain/web"
)

type LogoutController struct {
	LogoutUsecase web.LogoutUsecase
	Env           *bootstrap.Env
	Session       *session.Store
}

func (lc *LogoutController) Logout(c *fiber.Ctx) error {
	sessionId := c.Cookies("session_id")
	println(sessionId)

	err := lc.LogoutUsecase.ValidateSession(sessionId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&web.LogoutResponse{
			Message: err.Error(),
		})
	}

	content := lc.LogoutUsecase.DecryptSession(lc.Env.Key, sessionId)

	err = lc.LogoutUsecase.DeleteSession(c.Context(), content)
	if err != nil {
		panic(err)
	}

	c.Cookie(lc.LogoutUsecase.CreateExpiredCookie())

	return c.Status(fiber.StatusOK).JSON(web.LogoutResponse{
		Message: "Logout successful",
	})
}
