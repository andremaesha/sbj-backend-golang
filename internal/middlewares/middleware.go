package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"sbj-backend/domain"
	"sbj-backend/internal/logger"
	"strings"
	"time"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(&domain.ErrorResponse{
		Message: fiber.ErrUnauthorized.Message,
	})
}

func ErrorHandler(c *fiber.Ctx) error {
	defer func() error {
		if err := recover(); err != nil {
			requestId := fmt.Sprintf("request_%s", time.Now().Format("20060102150405"))

			buf := make([]byte, 2048) // Gunakan buffer yang lebih besar untuk stack trace yang lebih lengkap
			n := runtime.Stack(buf, true)
			stackTrace := strings.TrimSpace(string(buf[:n]))

			// Mencatat pesan panic dan stack trace
			logger.Error.Printf("RequestID: (%s) Recovered from panic: %v\nStack trace: %s", requestId, err, stackTrace)

			// Mengembalikan response error
			return c.Status(fiber.StatusInternalServerError).JSON(&domain.ErrorResponse{
				RequestId: requestId,
				Message:   fiber.ErrInternalServerError.Message,
			})
		}

		return nil
	}()

	return c.Next()
}
