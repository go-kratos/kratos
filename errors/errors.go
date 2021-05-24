package errors

import (
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/internal/httputil"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is describes the cause of the error with structured details.
// For more details see https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto.
type Error struct {
	Code    codes.Code `json:"code"`
	Message string     `json:"message"`
	// error deatils
	Domain   string            `json:"domain,omitempty"`
	Reason   string            `json:"reason,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: domain = %s reason = %s metadata = %v", e.Domain, e.Reason, e.Metadata)
}

// HTTPStatus return an HTTP error code.
func (e *Error) HTTPStatus() int {
	return httputil.StatusFromGRPCCode(e.Code)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(e.Code, e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Domain:   e.Domain,
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Domain == e.Domain && se.Reason == e.Reason
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := *e
	err.Metadata = md
	return &err
}

// New returns an error object for the code, message.
func New(code codes.Code, domain, reason, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Domain:  domain,
		Reason:  reason,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code codes.Code, domain, reason, format string, a ...interface{}) *Error {
	return New(code, domain, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code codes.Code, domain, reason, format string, a ...interface{}) error {
	return New(code, domain, reason, fmt.Sprintf(format, a...))
}

// Code returns the code for a particular error.
// It supports wrapped errors.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if se := FromError(err); err != nil {
		return se.Code
	}
	return codes.Unknown
}

// Domain returns the domain for a particular error.
// It supports wrapped errors.
func Domain(err error) string {
	if se := FromError(err); err != nil {
		return se.Domain
	}
	return ""
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if se := FromError(err); err != nil {
		return se.Reason
	}
	return ""
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				return New(
					gs.Code(),
					d.Domain,
					d.Reason,
					gs.Message(),
				).WithMetadata(d.Metadata)
			}
		}
	}
	return New(gs.Code(), "", "", err.Error())
}
