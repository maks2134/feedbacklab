// Package errors provides error wrapping and handling utilities.
package errors

import (
	"fmt"
)

// ErrorType represents the type of error.
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "not_found"
	ErrorTypeBadRequest    ErrorType = "bad_request"
	ErrorTypeInternal      ErrorType = "internal"
	ErrorTypeUnauthorized  ErrorType = "unauthorized"
	ErrorTypeForbidden     ErrorType = "forbidden"
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeAlreadyExists ErrorType = "already_exists"
)

// AppError represents an application error with type and context.
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
	Key     string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(key string, resource string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Err:     err,
		Key:     key,
	}
}

// NewBadRequestError creates a new bad request error.
func NewBadRequestError(key string, message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Err:     err,
		Key:     key,
	}
}

// NewInternalError creates a new internal error with wrapping.
func NewInternalError(layer string, operation string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: fmt.Sprintf("%s: %s failed", layer, operation),
		Err:     err,
		Key:     "common.error.internal",
	}
}

// NewValidationError creates a new validation error.
func NewValidationError(key string, message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
		Key:     key,
	}
}

// NewAlreadyExistsError creates a new already exists error.
func NewAlreadyExistsError(key string, resource string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeAlreadyExists,
		Message: fmt.Sprintf("%s already exists", resource),
		Err:     err,
		Key:     key,
	}
}
