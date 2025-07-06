package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/internal/encry"
	"time"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
	Session      *session.Store
}

func (lc *LoginController) Login(c *fiber.Ctx) error {
	request := new(domain.LoginRequest)
	sessionId := uuid.New().String()

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

	return c.Status(fiber.StatusOK).JSON(domain.LoginResponse{
		Id:      encryptSession,
		Email:   user.Email,
		Message: "success",
	})
}
