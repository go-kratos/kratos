package errors

import "fmt"

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

// WithMetadata with the metadata to error.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := *e
	err.Metadata = md
	return &err
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

// BadRequestError is Bad Request error that is mapped to a 400 response.
type BadRequestError struct {
	err *Error
}

func (e *BadRequestError) Unwrap() error { return e.err }
func (e *BadRequestError) Error() string { return e.err.Error() }

// NewBadRequest new a Bad Request error.
func NewBadRequest(domain, reason, message string) *BadRequestError {
	return &BadRequestError{err: New(400, domain, reason, message)}
}

// UnauthorizedError is Unauthorized error that is mapped to a 401 response.
type UnauthorizedError struct {
	err *Error
}

func (e *UnauthorizedError) Unwrap() error { return e.err }
func (e *UnauthorizedError) Error() string { return e.err.Error() }

// NewUnauthorized new a Unauthorized error.
func NewUnauthorized(domain, reason, message string) *UnauthorizedError {
	return &UnauthorizedError{err: New(401, domain, reason, message)}
}

// ForbiddenError is Forbidden error that is mapped to a 403 response.
type ForbiddenError struct {
	err *Error
}

func (e *ForbiddenError) Unwrap() error { return e.err }
func (e *ForbiddenError) Error() string { return e.err.Error() }

// NewForbidden new a Forbidden error.
func NewForbidden(domain, reason, message string) *ForbiddenError {
	return &ForbiddenError{err: New(403, domain, reason, message)}
}

// NotFoundError is Not Found error that is mapped to a 404 response.
type NotFoundError struct {
	err *Error
}

func (e *NotFoundError) Unwrap() error { return e.err }
func (e *NotFoundError) Error() string { return e.err.Error() }

// NewNotFound new a Not Found error.
func NewNotFound(domain, reason, message string) *NotFoundError {
	return &NotFoundError{err: New(403, domain, reason, message)}
}

// ConflictError Conflict error that is mapped to a 409 response.
type ConflictError struct {
	err *Error
}

func (e *ConflictError) Unwrap() error { return e.err }
func (e *ConflictError) Error() string { return e.err.Error() }

// NewConflict new a Conflict error.
func NewConflict(domain, reason, message string) *ConflictError {
	return &ConflictError{err: New(409, domain, reason, message)}
}

// InternalServerError is Internal Server error that is mapped to a 500 response.
type InternalServerError struct {
	err *Error
}

func (e *InternalServerError) Unwrap() error { return e.err }
func (e *InternalServerError) Error() string { return e.err.Error() }

// NewInternalServer new a Internal Server error.
func NewInternalServer(domain, reason, message string) *InternalServerError {
	return &InternalServerError{err: New(409, domain, reason, message)}
}

// ServiceUnavailableError is Service Unavailable response for the API, mapped to a HTTP 503 response.
type ServiceUnavailableError struct {
	err *Error
}

func (e *ServiceUnavailableError) Unwrap() error { return e.err }
func (e *ServiceUnavailableError) Error() string { return e.err.Error() }

// NewServiceUnavailable new a Service Unavailable error.
func NewServiceUnavailable(domain, reason, message string) *ServiceUnavailableError {
	return &ServiceUnavailableError{err: New(409, domain, reason, message)}
}
