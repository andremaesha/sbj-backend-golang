package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
)

type LogoutController struct {
	Env     *bootstrap.Env
	Session *session.Store
}

func (lc *LogoutController) Logout(c *fiber.Ctx) error {
	sess, err := lc.Session.Get(c)
	if err != nil {
		panic(err)
	}

	if err = sess.Destroy(); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(domain.LogoutResponse{
		Message: "Logout successful",
	})
}
