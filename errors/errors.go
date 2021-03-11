package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

var _ error = (*StatusError)(nil)

// StatusError contains an error response from the server.
type StatusError = Status

// Is matches each error in the chain with the target value.
func (e *StatusError) Is(target error) bool {
	err, ok := target.(*StatusError)
	if ok {
		return e.Code == err.Code
	}
	return false
}

// HTTPStatus returns the Status represented by se.
func (e *StatusError) HTTPStatus() int {
	switch e.Code {
	case 0:
		return http.StatusOK
	case 1:
		return http.StatusInternalServerError
	case 2:
		return http.StatusInternalServerError
	case 3:
		return http.StatusBadRequest
	case 4:
		return http.StatusRequestTimeout
	case 5:
		return http.StatusNotFound
	case 6:
		return http.StatusConflict
	case 7:
		return http.StatusForbidden
	case 8:
		return http.StatusTooManyRequests
	case 9:
		return http.StatusPreconditionFailed
	case 10:
		return http.StatusConflict
	case 11:
		return http.StatusBadRequest
	case 12:
		return http.StatusNotImplemented
	case 13:
		return http.StatusInternalServerError
	case 14:
		return http.StatusServiceUnavailable
	case 15:
		return http.StatusInternalServerError
	case 16:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s details = %+v", e.Code, e.Reason, e.Message, e.Details)
}

// Error returns a Status representing c and msg.
func Error(code int32, reason, message string) error {
	return &StatusError{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

// Errorf returns New(c, fmt.Sprintf(format, a...)).
func Errorf(code int32, reason, format string, a ...interface{}) error {
	return Error(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the status code.
func Code(err error) int32 {
	if err == nil {
		return 0 // ok
	}
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code
	}
	return 2 // unknown
}

// Reason returns the status for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Reason
	}
	return UnknownReason
}

// FromError returns status error.
func FromError(err error) (*StatusError, bool) {
	if se := new(StatusError); errors.As(err, &se) {
		return se, true
	}
	return nil, false
}
