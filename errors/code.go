package errors

import (
	"net/http"
)

// BadRequest generates a 400 error.
func BadRequest(errs ...ErrorInfo) error {
	return &Error{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Details: []interface{}{errs},
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(errs ...ErrorInfo) error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Details: []interface{}{errs},
	}
}
