package errors

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
)

const (
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

var _ error = &StatusError{}

// StatusError contains an error response from the server.
type StatusError struct {
	// Code is the gRPC response status code and will always be populated.
	Code int32 `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Details provide more context to an error.
	Details []proto.Message `json:"details"`
}

// WithDetails provided details messages appended to the errors.
func (e *StatusError) WithDetails(details ...proto.Message) {
	e.Details = append(e.Details, details...)
}

// Is matches each error in the chain with the target value.
func (e *StatusError) Is(target error) bool {
	err, ok := target.(*StatusError)
	if ok {
		return e.Code == err.Code
	}
	return false
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", e.Code, e.Message, e.Details)
}

// Error returns a Status representing c and msg.
func Error(code int32, message string, details ...proto.Message) error {
	return &StatusError{Code: code, Message: message, Details: details}
}

// Errorf returns New(c, fmt.Sprintf(format, a...)).
func Errorf(code int32, format string, a ...interface{}) error {
	return Error(code, fmt.Sprintf(format, a...))
}

// Reason returns the gRPC status for a particular error.
// It supports wrapped errors.
func Reason(err error) *ErrorItem {
	if se := new(StatusError); errors.As(err, &se) {
		for _, d := range se.Details {
			if e, ok := d.(*ErrorItem); ok {
				return e
			}
		}
	}
	return &ErrorItem{Reason: UnknownReason}
}
