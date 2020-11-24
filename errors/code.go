package errors

import (
	"net/http"
)

// BadRequest generates a 400 error.
func BadRequest(err ErrorInfo) error {
	return &Error{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Details: []interface{}{err},
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(err ErrorInfo) error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Details: []interface{}{err},
	}
}
