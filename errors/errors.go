package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error contains an error response from the server.
// For more details see https://github.com/go-kratos/kratos/issues/858.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d message = %s", e.Code, e.Message)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if target := new(Error); errors.As(err, &target) {
		return target.Code == e.Code
	}
	return false
}

// New returns an error object for the code, message.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, domain, reason, format string, a ...interface{}) *ErrorInfo {
	return &ErrorInfo{
		err:    Newf(code, format, a...),
		Domain: domain,
		Reason: reason,
	}
}

// Code returns the code for a particular error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if target := new(Error); errors.As(err, &target) {
		return target.Code
	}
	return http.StatusInternalServerError
}

// Domain returns the domain for a particular error.
// It supports wrapped errors.
func Domain(err error) string {
	if target := new(ErrorInfo); errors.As(err, &target) {
		return target.Domain
	}
	return ""
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if target := new(ErrorInfo); errors.As(err, &target) {
		return target.Reason
	}
	return ""
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if target := new(Error); errors.As(err, &target) {
		return target
	}
	return New(http.StatusInternalServerError, err.Error())
}
