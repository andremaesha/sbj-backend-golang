package controller

import (
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

type LoginController struct {
	LoginUsecase web.LoginUsecase
	Env          *bootstrap.Env
	Session      *session.Store
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	request := new(web.LoginRequest)
	sessionId := uuid.New().String()

	if c.BodyParser(request) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "error with you're json body"})
	}

	// Validate request
	if err := validator.ValidateStruct(request); err != nil {
		return validator.HandleValidationErrors(c, err)
	}

	user, err := lc.LoginUsecase.GetUserByEmail(c.Context(), request.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "user not found with the given email"})
	}

	err = lc.LoginUsecase.ValidateUserCredentials(user.Password, request.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	err = lc.LoginUsecase.ValidateUserVerified(user.Verified)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	err = lc.LoginUsecase.SetSession(c.Context(), lc.Env.RedisExpireTime, sessionId, user)
	if err != nil {
		panic(err)
	}

	encryptSession := lc.LoginUsecase.EncryptSession(lc.Env.Key, sessionId)

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    encryptSession,
		Expires:  time.Now().Add(time.Minute * time.Duration(lc.Env.RedisExpireTime)),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(web.LoginResponse{
		Email:   user.Email,
		Message: "success",
	})
}
