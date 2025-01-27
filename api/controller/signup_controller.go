package controller

import (
	"github.com/gofiber/fiber/v2"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/encry"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SignupController) Signup(c *fiber.Ctx) error {
	request := new(domain.SignupRequest)

	if c.BodyParser(request) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "error with you're json body"})
	}

	_, err := sc.SignupUsecase.GetUserByEmail(c.Context(), request.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(domain.ErrorResponse{Message: "user already exists with the given email"})
	}

	encryptedPassword, err := encry.HashPassword(request.Password)
	if err != nil {
		panic(err)
	}

	user := &domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  encryptedPassword,
	}

	err = sc.SignupUsecase.Create(c.Context(), user)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(domain.SignupResponse{
		Message: "success",
	})
}
