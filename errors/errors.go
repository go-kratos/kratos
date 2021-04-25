package errors

import (
	"errors"
	"fmt"
	"strings"
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

// WithMetadata returns an MD formed by the mapping of key, value ...
// WithMetadata panics if len(kv) is odd.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
func (e *Error) WithMetadata(kv ...string) *Error {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := make(map[string]string, len(kv)/2)
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = strings.ToLower(s)
			continue
		}
		md[key] = s
	}
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
	return fmt.Sprintf("error: code = %d message = %s reason = %s domain = %s", e.Code, e.Message, e.Reason, e.Domain)
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
	if target := new(Error); errors.As(err, &target) {
		return target.Code
	}
	return 500
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if target := new(Error); errors.As(err, &target) {
		return target.Reason
	}
	return ""
}
