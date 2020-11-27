package status

import (
	"net/http"
)

// BadRequest generates a 400 error.
func BadRequest(err *ErrorInfo) error {
	return &Status{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Details: []interface{}{err},
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(err *ErrorInfo) error {
	return &Status{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Details: []interface{}{err},
	}
}
