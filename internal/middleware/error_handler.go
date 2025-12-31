package middleware

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/service"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

// ErrorDetails contains error information
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorHandler is a centralized error handling middleware
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			// Handle echo.HTTPError
			var httpError *echo.HTTPError
			if errors.As(err, &httpError) {
				return handleHTTPError(c, httpError)
			}

			// Handle application errors
			var appError service.AppError
			if errors.As(err, &appError) {
				return handleAppError(c, appError)
			}

			// Handle validation errors
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return handleValidationError(c, validationErrors)
			}

			// Handle unknown errors
			c.Logger().Error("Unhandled error:", err)
			appErr := service.NewInternalServerError("An unexpected error occurred")
			return handleAppError(c, appErr)
		}
	}
}

// handleAppError handles structured application errors
func handleAppError(c echo.Context, appErr service.AppError) error {
	response := ErrorResponse{
		Error: ErrorDetails{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	c.Logger().Errorf("Application error: %v", appErr)
	return c.JSON(appErr.HTTPStatus, response)
}

// handleHTTPError handles echo HTTP errors
func handleHTTPError(c echo.Context, httpErr *echo.HTTPError) error {
	code := "HTTP_ERROR"
	message := "HTTP error occurred"

	switch httpErr.Code {
	case http.StatusNotFound:
		code = service.ErrCodeNotFound
		message = "Resource not found"
	case http.StatusUnauthorized:
		code = service.ErrCodeUnauthorized
		message = "Unauthorized"
	case http.StatusForbidden:
		code = service.ErrCodeForbidden
		message = "Forbidden"
	case http.StatusBadRequest:
		code = service.ErrCodeBadRequest
		message = "Bad request"
	}

	response := ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
			Details: httpErr.Message.(string),
		},
	}

	return c.JSON(httpErr.Code, response)
}

// handleValidationError handles validation errors
func handleValidationError(c echo.Context, validationErrors validator.ValidationErrors) error {
	details := "Validation failed for the following fields:"
	for _, err := range validationErrors {
		details += " " + err.Field() + " (" + err.Tag() + ")"
	}

	response := ErrorResponse{
		Error: ErrorDetails{
			Code:    service.ErrCodeValidation,
			Message: "Validation failed",
			Details: details,
		},
	}

	return c.JSON(http.StatusBadRequest, response)
}