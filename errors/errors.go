package errors

import "fmt"

const (
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Domain  string `json:"domain"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d message = %s reason = %s domain = %s", e.Code, e.Message, e.Reason, e.Domain)
}

func New(code int, message, reason, domain string) *Error {
	return &Error{
		Code:    code,
		Domain:  domain,
		Reason:  reason,
		Message: message,
	}
}

type BadRequestError struct {
	err *Error
}

func (e *BadRequestError) Unwrap() error { return e.err }
func (e *BadRequestError) Error() string { return e.err.Error() }

func NewBadRequest(message, reason, domain string) *BadRequestError {
	return &BadRequestError{err: New(400, domain, reason, message)}
}
