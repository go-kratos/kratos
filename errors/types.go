package errors

// BadRequest new BadRequest error that is mapped to a 400 response.
func BadRequest(domain, reason, message string) *Error {
	return New(400, domain, reason, message)
}

// IsBadRequest determines if err is an error which indicates an BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return Code(err) == 400
}

// Unauthorized new Unauthorized error that is mapped to a 401 response.
func Unauthorized(domain, reason, message string) *Error {
	return New(401, domain, reason, message)
}

// IsUnauthorized determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return Code(err) == 401
}

// Forbidden new Forbidden error that is mapped to a 403 response.
func Forbidden(domain, reason, message string) *Error {
	return New(403, domain, reason, message)
}

// IsForbidden determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == 403
}

// NotFound new NotFound error that is mapped to a 404 response.
func NotFound(domain, reason, message string) *Error {
	return New(404, domain, reason, message)
}

// IsNotFound determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == 404
}

// Conflict new Conflict error that is mapped to a 409 response.
func Conflict(domain, reason, message string) *Error {
	return New(409, domain, reason, message)
}

// IsConflict determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsConflict(err error) bool {
	return Code(err) == 409
}

// InternalServer new InternalServer error that is mapped to a 500 response.
func InternalServer(domain, reason, message string) *Error {
	return New(500, domain, reason, message)
}

// IsInternalServer determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsInternalServer(err error) bool {
	return Code(err) == 500
}

// ServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ServiceUnavailable(domain, reason, message string) *Error {
	return New(503, domain, reason, message)
}

// IsServiceUnavailable determines if err is an error which indicates an Unauthorized error.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}
