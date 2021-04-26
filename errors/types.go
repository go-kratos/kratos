package errors

import (
	"errors"
	"fmt"
)

// ErrorInfo is describes the cause of the error with structured details.
// For more details see https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto.
type ErrorInfo struct {
	err      *Error
	Domain   string            `json:"domain"`
	Reason   string            `json:"reason"`
	Metadata map[string]string `json:"metadata"`
}

func (e *ErrorInfo) Error() string {
	return fmt.Sprintf("error: domain = %s reason = %s", e.Domain, e.Reason)
}
func (e *ErrorInfo) Unwrap() error {
	return e.err
}

// Is matches each error in the chain with the target value.
func (e *ErrorInfo) Is(err error) bool {
	if target := new(ErrorInfo); errors.As(err, &target) {
		return target.Domain == e.Domain && target.Reason == e.Reason
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *ErrorInfo) WithMetadata(md map[string]string) *ErrorInfo {
	err := *e
	err.Metadata = md
	return &err
}
