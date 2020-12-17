package http

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
)

var errMapping = map[int32]int{
	0:  http.StatusOK,
	1:  http.StatusInternalServerError,
	2:  http.StatusInternalServerError,
	3:  http.StatusBadRequest,
	4:  http.StatusRequestTimeout,
	5:  http.StatusNotFound,
	6:  http.StatusConflict,
	7:  http.StatusForbidden,
	9:  http.StatusPreconditionFailed,
	10: http.StatusConflict,
	11: http.StatusBadRequest,
	12: http.StatusNotImplemented,
	13: http.StatusInternalServerError,
	14: http.StatusServiceUnavailable,
	15: http.StatusInternalServerError,
	16: http.StatusUnauthorized,
}

// StatusError converts error to http error.
func StatusError(err error) (*errors.StatusError, int) {
	se, ok := err.(*errors.StatusError)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Message: "Unknown: " + err.Error(),
		}
	}
	if code, ok := errMapping[se.Code]; ok {
		return se, code
	}
	return se, http.StatusInternalServerError
}
