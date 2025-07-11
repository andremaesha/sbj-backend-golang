package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
)

type ProductsController struct {
	ProductsUsecase web.ProductsUsecase
	Env             *bootstrap.Env
	Session         *session.Store
}

func (p *ProductsController) Product(c *fiber.Ctx) error {
	id := c.Query("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Message: "id parameter is required",
		})
	}

	response, err := p.ProductsUsecase.Product(c.Context(), id)
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
