package middlewares

import (
	"fmt"
	"sbj-backend/internal/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ResponseLogger(c *fiber.Ctx) error {
	logger.CheckAndRotateLog()

	urlDetail := fmt.Sprintf("URL DETAIL: %s", c.BaseURL()+c.OriginalURL())

	start := time.Now()

	// Log request details
	reqDetail := fmt.Sprintf("Started %s %s for %s", c.Method(), c.OriginalURL(), c.IP())

	// Process request
	err := c.Next()

	// Log response details
	duration := time.Since(start)
	resDetail := fmt.Sprintf("Completed %s in %v with status %d", c.Path(), duration, c.Response().StatusCode())

	requestBody := string(c.Body())
	reqBody := fmt.Sprintf("Request Body: %s", requestBody)

	// Log response body
	responseBody := string(c.Response().Body())
	resBody := fmt.Sprintf("Response Body: %s", responseBody)

	logger.Info.Println(urlDetail, reqDetail, resDetail, reqBody, resBody)

	return err
}
