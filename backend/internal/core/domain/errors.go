package domain

import (
	"errors"
	"fmt"
)

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrWeakPassword      = errors.New("password does not meet requirements")

	// URL errors
	ErrShortURLNotFound    = errors.New("short URL not found")
	ErrURLNotFound         = errors.New("URL not found")
	ErrShortCodeExists     = errors.New("short code already exists")
	ErrInvalidURL          = errors.New("invalid URL format")
	ErrURLExpired          = errors.New("URL has expired")
	ErrURLInactive         = errors.New("URL is inactive")
	ErrCustomAliasInvalid  = errors.New("custom alias is invalid")
	ErrCustomAliasTooLong  = errors.New("custom alias is too long")

	// Authentication errors
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("token has expired")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidPassword    = errors.New("invalid password")

	// Validation errors
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrMissingField        = errors.New("required field is missing")
	ErrInvalidFieldFormat  = errors.New("invalid field format")
	ErrInvalidShortCode    = errors.New("invalid short code")

	// Database errors
	ErrDatabaseConnection  = errors.New("database connection failed")
	ErrDatabaseQuery       = errors.New("database query failed")
	ErrDatabaseTransaction = errors.New("database transaction failed")

	// Rate limiting errors
	ErrRateLimitExceeded   = errors.New("rate limit exceeded")
	ErrTooManyRequests     = errors.New("too many requests")

	// External service errors
	ErrExternalService     = errors.New("external service error")
	ErrGeolocationService  = errors.New("geolocation service error")
)

type DomainError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func NewDomainError(errorType, message string, code int) *DomainError {
	return &DomainError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}

func NewValidationError(field, message string) *DomainError {
	return &DomainError{
		Type:    "validation_error",
		Message: fmt.Sprintf("%s: %s", field, message),
		Code:    400,
	}
}

func NewNotFoundError(resource string) *DomainError {
	return &DomainError{
		Type:    "not_found",
		Message: fmt.Sprintf("%s not found", resource),
		Code:    404,
	}
}

func NewConflictError(resource string) *DomainError {
	return &DomainError{
		Type:    "conflict",
		Message: fmt.Sprintf("%s already exists", resource),
		Code:    409,
	}
}

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:    "unauthorized",
		Message: message,
		Code:    401,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Type:    "forbidden",
		Message: message,
		Code:    403,
	}
}

func NewInternalError(message string) *DomainError {
	return &DomainError{
		Type:    "internal_error",
		Message: message,
		Code:    500,
	}
}