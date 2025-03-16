package service

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
}

func (e AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error codes
const (
	ErrCodeNotFound         = "RESOURCE_NOT_FOUND"
	ErrCodeConflict         = "RESOURCE_CONFLICT"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeInternalServer   = "INTERNAL_SERVER_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeBadRequest       = "BAD_REQUEST"
)

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, identifier string) AppError {
	return AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		Details:    fmt.Sprintf("No %s found with identifier: %s", resource, identifier),
		HTTPStatus: http.StatusNotFound,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(resource, field, value string) AppError {
	return AppError{
		Code:       ErrCodeConflict,
		Message:    fmt.Sprintf("%s already exists", resource),
		Details:    fmt.Sprintf("A %s with %s '%s' already exists", resource, field, value),
		HTTPStatus: http.StatusConflict,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string) AppError {
	return AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		Details:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string) AppError {
	return AppError{
		Code:       ErrCodeInternalServer,
		Message:    "Internal server error",
		Details:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) AppError {
	return AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Unauthorized",
		Details:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) AppError {
	return AppError{
		Code:       ErrCodeForbidden,
		Message:    "Forbidden",
		Details:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) AppError {
	return AppError{
		Code:       ErrCodeBadRequest,
		Message:    "Bad request",
		Details:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}