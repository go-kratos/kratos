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
	Code     int               `json:"code"`
	Domain   string            `json:"domain"`
	Reason   string            `json:"reason"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata"`
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := *e
	err.Metadata = md
	return &err
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if target := new(Error); errors.As(err, &target) {
		return target.Code == e.Code && target.Reason == e.Reason
	}
	return false
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d domain = %s reason = %s message = %s", e.Code, e.Domain, e.Reason, e.Message)
}

// New returns an error object for the code, message and error info.
func New(code int, domain, reason, message string) *Error {
	return &Error{
		Code:    code,
		Domain:  domain,
		Reason:  reason,
		Message: message,
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
	if target := new(Error); errors.As(err, &target) {
		return target.Domain
	}
	return ""
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if target := new(Error); errors.As(err, &target) {
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
	return New(http.StatusInternalServerError, "", "", err.Error())
}
