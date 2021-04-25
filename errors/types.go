package errors

import "net/http"

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(domain, reason, message string) *Error {
	return New(http.StatusBadRequest, domain, reason, message)
}

// IsBadRequest determines if err is an error which indicates a BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return Code(err) == http.StatusBadRequest
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(domain, reason, message string) *Error {
	return New(http.StatusUnauthorized, domain, reason, message)
}

// IsUnauthorized determines if err is an error which indicates a Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return Code(err) == http.StatusUnauthorized
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(domain, reason, message string) *Error {
	return New(http.StatusForbidden, domain, reason, message)
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == http.StatusForbidden
}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(domain, reason, message string) *Error {
	return New(http.StatusNotFound, domain, reason, message)
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == http.StatusNotFound
}

// Conflict new Conflict error that is mapped to a 409 response.
func Conflict(domain, reason, message string) *Error {
	return New(http.StatusConflict, domain, reason, message)
}

// IsConflict determines if err is an error which indicates a Conflict error.
// It supports wrapped errors.
func IsConflict(err error) bool {
	return Code(err) == http.StatusConflict
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func InternalServer(domain, reason, message string) *Error {
	return New(http.StatusInternalServerError, domain, reason, message)
}

// IsInternalServer determines if err is an error which indicates an InternalServer error.
// It supports wrapped errors.
func IsInternalServer(err error) bool {
	return Code(err) == http.StatusInternalServerError
}

// ServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ServiceUnavailable(domain, reason, message string) *Error {
	return New(http.StatusServiceUnavailable, domain, reason, message)
}

// IsServiceUnavailable determines if err is an error which indicates a ServiceUnavailable error.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
	return Code(err) == http.StatusServiceUnavailable
}
