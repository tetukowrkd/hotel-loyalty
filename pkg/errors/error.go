package errors

import "net/http"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
	}

	ErrInternal = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	}

	ErrEmailExists = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Email Already Exists",
	}
)
