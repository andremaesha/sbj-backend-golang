package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
	"sbj-backend/internal/validator"
)

type ProductsController struct {
	ProductsUsecase web.ProductsUsecase
	Env             *bootstrap.Env
	Session         *session.Store
}

func (p *ProductsController) Product(c *fiber.Ctx) error {
	request := new(web.ProductsRequest)
	request.Id = c.Query("id")

	// Validate request
	if err := validator.ValidateStruct(request); err != nil {
		return validator.HandleValidationErrors(c, err)
	}

	response, err := p.ProductsUsecase.Product(c.Context(), request.Id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (p *ProductsController) Products(c *fiber.Ctx) error {
	response, err := p.ProductsUsecase.Products(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (p *ProductsController) CreateProduct(c *fiber.Ctx) error {
	sessionId := c.Cookies("session_id")

	err := p.ProductsUsecase.ValidatePermission(c.Context(), p.Env.Key, sessionId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	request := new(web.ProductRequest)

	if err = c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: "error with your json body"})
	}

	// Validate request
	if err = validator.ValidateStruct(request); err != nil {
		return validator.HandleValidationErrors(c, err)
	}

	err = p.ProductsUsecase.ProductCreate(c.Context(), p.Env, request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(web.ProductResponse{
		ResponseMessage: "Product created successfully",
	})
}
