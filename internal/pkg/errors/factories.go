package errors

import "net/http"

func NewBadRequestError(msg string, err ...error) *AppError {
	return newAppError("BAD_REQUEST", msg, http.StatusBadRequest, err...)
}

func NewUnauthorizedError(msg string, err ...error) *AppError {
	return newAppError("UNAUTHORIZED", msg, http.StatusUnauthorized, err...)
}

func NewNotFoundError(msg string, err ...error) *AppError {
	return newAppError("NOT_FOUND", msg, http.StatusNotFound, err...)
}

func NewConflictError(msg string, err ...error) *AppError {
	return newAppError("CONFLICT", msg, http.StatusConflict, err...)
}

func NewInternalError(err ...error) *AppError {
	return newAppError("INTERNAL_ERROR", "An unexpected error occurred", http.StatusInternalServerError, err...)
}

func NewForbiddenError() *AppError {
	return newAppError("FORBIDDEN", "FORBIDDEN", http.StatusForbidden)
}

func newAppError(code, msg string, statuCode int, err ...error) *AppError {
	var underlyingErr error
	if len(err) > 0 {
		underlyingErr = err[0]
	}

	return &AppError{
		StatusCode: statuCode,
		Code:       code,
		Message:    msg,
		Err:        underlyingErr,
	}
}
