package validator

import (
	"reflect"
	"sbj-backend/domain"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

// Initialize initializes the validator
func Initialize() {
	validate = validator.New()

	// Register function to get tag name from json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		Initialize()
	}
	return validate
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	if validate == nil {
		Initialize()
	}
	return validate.Struct(s)
}

// HandleValidationErrors handles validation errors and returns a fiber response
func HandleValidationErrors(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Check if the error is a validation error
	if _, ok := err.(validator.ValidationErrors); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get validation errors
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]string)

	// Format validation errors
	for _, e := range validationErrors {
		errorMessages[e.Field()] = formatErrorMessage(e)
	}

	return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
		Message: "Validation failed",
		Errors:  errorMessages,
	})
}

// formatErrorMessage formats a validation error message
func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value must be greater than or equal to " + e.Param()
	case "max":
		return "Value must be less than or equal to " + e.Param()
	case "oneof":
		return "Value must be one of " + e.Param()
	default:
		return "Invalid value"
	}
}
