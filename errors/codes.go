package errors

// Cancelled The operation was cancelled, typically by the caller.
// HTTP Mapping: 499 Client Closed Request
func Cancelled(reason, message string) error {
	return &StatusError{
		Code:    1,
		Message: "Canceled",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Unknown error.
// HTTP Mapping: 500 Internal Server Error
func Unknown(reason, message string) error {
	return &StatusError{
		Code:    2,
		Message: "Unknown",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// InvalidArgument The client specified an invalid argument.
// HTTP Mapping: 400 Bad Request
func InvalidArgument(reason, message string) error {
	return &StatusError{
		Code:    3,
		Message: "InvalidArgument",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// DeadlineExceeded The deadline expired before the operation could complete.
// HTTP Mapping: 504 Gateway Timeout
func DeadlineExceeded(reason, message string) error {
	return &StatusError{
		Code:    4,
		Message: "DeadlineExceeded",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// NotFound Some requested entity (e.g., file or directory) was not found.
// HTTP Mapping: 404 Not Found
func NotFound(reason, message string) error {
	return &StatusError{
		Code:    5,
		Message: "NotFound",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// AlreadyExists The entity that a client attempted to create (e.g., file or directory) already exists.
// HTTP Mapping: 409 Conflict
func AlreadyExists(reason, message string) error {
	return &StatusError{
		Code:    6,
		Message: "AlreadyExists",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// PermissionDenied The caller does not have permission to execute the specified operation.
// HTTP Mapping: 403 Forbidden
func PermissionDenied(reason, message string) error {
	return &StatusError{
		Code:    7,
		Message: "PermissionDenied",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// ResourceExhausted Some resource has been exhausted, perhaps a per-user quota, or
// perhaps the entire file system is out of space.
// HTTP Mapping: 429 Too Many Requests
func ResourceExhausted(reason, message string) error {
	return &StatusError{
		Code:    8,
		Message: "ResourceExhausted",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// FailedPrecondition The operation was rejected because the system is not in a state
// required for the operation's execution.
// HTTP Mapping: 400 Bad Request
func FailedPrecondition(reason, message string) error {
	return &StatusError{
		Code:    9,
		Message: "FailedPrecondition",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Aborted The operation was aborted, typically due to a concurrency issue such as
// a sequencer check failure or transaction abort.
// HTTP Mapping: 409 Conflict
func Aborted(reason, message string) error {
	return &StatusError{
		Code:    10,
		Message: "Aborted",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// OutOfRange The operation was attempted past the valid range.  E.g., seeking or
// reading past end-of-file.
// HTTP Mapping: 400 Bad Request
func OutOfRange(reason, message string) error {
	return &StatusError{
		Code:    11,
		Message: "OutOfRange",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Unimplemented The operation is not implemented or is not supported/enabled in this service.
// HTTP Mapping: 501 Not Implemented
func Unimplemented(reason, message string) error {
	return &StatusError{
		Code:    12,
		Message: "Unimplemented",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Internal This means that some invariants expected by the
// underlying system have been broken.  This error code is reserved
// for serious errors.
//
// HTTP Mapping: 500 Internal Server Error
func Internal(reason, message string) error {
	return &StatusError{
		Code:    13,
		Message: "Internal",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Unavailable The service is currently unavailable.
// HTTP Mapping: 503 Service Unavailable
func Unavailable(reason, message string) error {
	return &StatusError{
		Code:    14,
		Message: "Unavailable",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// DataLoss Unrecoverable data loss or corruption.
// HTTP Mapping: 500 Internal Server Error
func DataLoss(reason, message string) error {
	return &StatusError{
		Code:    15,
		Message: "DataLoss",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}

// Unauthorized The request does not have valid authentication credentials for the operation.
// HTTP Mapping: 401 Unauthorized
func Unauthorized(reason, message string) error {
	return &StatusError{
		Code:    16,
		Message: "Unauthenticated",
		Details: []interface{}{
			&ErrorInfo{Reason: reason, Message: message},
		},
	}
}
