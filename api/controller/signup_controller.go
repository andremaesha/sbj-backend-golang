package controller

import (
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SignupController struct {
	SignupUsecase web.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SignupController) Signup(c *fiber.Ctx) error {
	request := new(web.SignupRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "error with your json body"})
	}

	// Validate request
	if err := validator.ValidateStruct(request); err != nil {
		return validator.HandleValidationErrors(c, err)
	}

	_, err := sc.SignupUsecase.GetUserByEmail(c.Context(), request.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(domain.ErrorResponse{Message: "user already exists with the given email"})
	}

	err = sc.SignupUsecase.CreateUser(sc.Env, c.Context(), request)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusCreated).JSON(web.SignupResponse{
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

	return c.Status(fiber.StatusOK).JSON(web.UploadAvatarResponse{
		Id:        idFile,
		Filename:  newFileName,
		UrlAvatar: responseCloudinary.SecureUrl,
		Message:   "success",
	})
}
