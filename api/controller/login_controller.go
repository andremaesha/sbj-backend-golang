package controller

import (
	"github.com/gofiber/fiber/v2"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/encry"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	request := new(domain.LoginRequest)

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

	return c.Status(fiber.StatusOK).JSON(domain.LoginResponse{
		Email:   user.Email,
		Message: "success",
	})
}
