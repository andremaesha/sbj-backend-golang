package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	return c.Status(fiber.StatusCreated).JSON(domain.SignupResponse{
		Message: "success",
	})
}

func (sc *SignupController) UploadAvatar(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "file is required"})
	}

	idFile := uuid.New().String()

	newFileName := idFile + "_ori_" + file.Filename

	responseCloudinary, err := sc.SignupUsecase.UploadAvatar(sc.Env, file)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(domain.UploadAvatarResponse{
		Id:        idFile,
		Filename:  newFileName,
		UrlAvatar: responseCloudinary.SecureUrl,
		Message:   "success",
	})
}
