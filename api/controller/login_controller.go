package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/encry"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
	Session      *session.Store
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	request := new(domain.LoginRequest)
	idSession := uuid.New()

	if c.BodyParser(request) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "error with you're json body"})
	}

	user, err := lc.LoginUsecase.GetUserByEmail(c.Context(), request.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "user not found with the given email"})
	}

	if !encry.VerifyPassword(user.Password, request.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "invalid credentials"})
	}

	sess, err := lc.Session.Get(c)
	if err != nil {
		panic(err)
	}

	sess.Set("user", user.Email)
	sess.Set("userid", idSession.String())
	if err = sess.Save(); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(domain.LoginResponse{
		Id:      "",
		Email:   user.Email,
		Message: "success",
	})
}
