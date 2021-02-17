package http

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	// References: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
	codesMapping = map[int32]int{
		0:  http.StatusOK,
		1:  http.StatusInternalServerError,
		2:  http.StatusInternalServerError,
		3:  http.StatusBadRequest,
		4:  http.StatusRequestTimeout,
		5:  http.StatusNotFound,
		6:  http.StatusConflict,
		7:  http.StatusForbidden,
		8:  http.StatusTooManyRequests,
		9:  http.StatusPreconditionFailed,
		10: http.StatusConflict,
		11: http.StatusBadRequest,
		12: http.StatusNotImplemented,
		13: http.StatusInternalServerError,
		14: http.StatusServiceUnavailable,
		15: http.StatusInternalServerError,
		16: http.StatusUnauthorized,
	}
	statusMapping = map[int]int32{
		http.StatusOK:                  0,
		http.StatusBadRequest:          3,
		http.StatusRequestTimeout:      4,
		http.StatusNotFound:            5,
		http.StatusConflict:            6,
		http.StatusForbidden:           7,
		http.StatusUnauthorized:        16,
		http.StatusPreconditionFailed:  9,
		http.StatusNotImplemented:      12,
		http.StatusInternalServerError: 13,
		http.StatusServiceUnavailable:  14,
	}
)

// StatusError converts error to status error.
func StatusError(err error) (*errors.StatusError, int) {
	se, ok := errors.FromError(err)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Reason:  "Unknown",
			Message: "Unknown: " + err.Error(),
		}
	}
	if status, ok := codesMapping[se.Code]; ok {
		return se, status
	}
	return se, http.StatusInternalServerError
}
