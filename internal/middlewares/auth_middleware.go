package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"sbj-backend/bootstrap"
	"sbj-backend/domain"
	"sbj-backend/domain/web"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(env *bootstrap.Env, session *session.Store, authUsecase web.AuthUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get session cookie
		sessionCookie := c.Cookies("session_id")
		if sessionCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Message: "Unauthorized: No session found",
			})
		}

		// Decrypt session ID using auth usecase
		sessionID, err := authUsecase.DecryptSessionID(env.Key, sessionCookie)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Message: "Unauthorized: Invalid session",
			})
		}

		// Get user from session using auth usecase
		user, err := authUsecase.GetUserFromSession(c.Context(), sessionID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Message: "Unauthorized: Session expired or invalid",
			})
		}

		// Store user in context for later use
		c.Locals("user", user)

		return c.Next()
	}
}

// AdminRoleMiddleware checks if the authenticated user has admin role
func AdminRoleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by AuthMiddleware)
		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Message: "Unauthorized: Authentication required",
			})
		}

		// Check if user has admin role
		if user.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(domain.ErrorResponse{
				Message: "Forbidden: Admin role required",
			})
		}

		return c.Next()
	}
}
