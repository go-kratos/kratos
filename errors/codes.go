package errors

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Cancelled The operation was cancelled, typically by the caller.
// HTTP Mapping: 499 Client Closed Request
func Cancelled(reason, format string, a ...string) error {
	return &StatusError{
		Code:    1,
		Message: "Canceled",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Unknown error.
// HTTP Mapping: 500 Internal Server Error
func Unknown(reason, format string, a ...string) error {
	return &StatusError{
		Code:    2,
		Message: "Unknown",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// InvalidArgument The client specified an invalid argument.
// HTTP Mapping: 400 Bad Request
func InvalidArgument(reason, format string, a ...string) error {
	return &StatusError{
		Code:    3,
		Message: "InvalidArgument",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// DeadlineExceeded The deadline expired before the operation could complete.
// HTTP Mapping: 504 Gateway Timeout
func DeadlineExceeded(reason, format string, a ...string) error {
	return &StatusError{
		Code:    4,
		Message: "DeadlineExceeded",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// NotFound Some requested entity (e.g., file or directory) was not found.
// HTTP Mapping: 404 Not Found
func NotFound(reason, format string, a ...string) error {
	return &StatusError{
		Code:    5,
		Message: "NotFound",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// AlreadyExists The entity that a client attempted to create (e.g., file or directory) already exists.
// HTTP Mapping: 409 Conflict
func AlreadyExists(reason, format string, a ...string) error {
	return &StatusError{
		Code:    6,
		Message: "AlreadyExists",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// PermissionDenied The caller does not have permission to execute the specified operation.
// HTTP Mapping: 403 Forbidden
func PermissionDenied(reason, format string, a ...string) error {
	return &StatusError{
		Code:    7,
		Message: "PermissionDenied",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// ResourceExhausted Some resource has been exhausted, perhaps a per-user quota, or
// perhaps the entire file system is out of space.
// HTTP Mapping: 429 Too Many Requests
func ResourceExhausted(reason, format string, a ...string) error {
	return &StatusError{
		Code:    8,
		Message: "ResourceExhausted",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// FailedPrecondition The operation was rejected because the system is not in a state
// required for the operation's execution.
// HTTP Mapping: 400 Bad Request
func FailedPrecondition(reason, format string, a ...string) error {
	return &StatusError{
		Code:    9,
		Message: "FailedPrecondition",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Aborted The operation was aborted, typically due to a concurrency issue such as
// a sequencer check failure or transaction abort.
// HTTP Mapping: 409 Conflict
func Aborted(reason, format string, a ...string) error {
	return &StatusError{
		Code:    10,
		Message: "Aborted",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// OutOfRange The operation was attempted past the valid range.  E.g., seeking or
// reading past end-of-file.
// HTTP Mapping: 400 Bad Request
func OutOfRange(reason, format string, a ...string) error {
	return &StatusError{
		Code:    11,
		Message: "OutOfRange",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Unimplemented The operation is not implemented or is not supported/enabled in this service.
// HTTP Mapping: 501 Not Implemented
func Unimplemented(reason, format string, a ...string) error {
	return &StatusError{
		Code:    12,
		Message: "Unimplemented",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Internal This means that some invariants expected by the
// underlying system have been broken.  This error code is reserved
// for serious errors.
//
// HTTP Mapping: 500 Internal Server Error
func Internal(reason, format string, a ...string) error {
	return &StatusError{
		Code:    13,
		Message: "Internal",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Unavailable The service is currently unavailable.
// HTTP Mapping: 503 Service Unavailable
func Unavailable(reason, format string, a ...string) error {
	return &StatusError{
		Code:    14,
		Message: "Unavailable",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// DataLoss Unrecoverable data loss or corruption.
// HTTP Mapping: 500 Internal Server Error
func DataLoss(reason, format string, a ...string) error {
	return &StatusError{
		Code:    15,
		Message: "DataLoss",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}

// Unauthorized The request does not have valid authentication credentials for the operation.
// HTTP Mapping: 401 Unauthorized
func Unauthorized(reason, format string, a ...string) error {
	return &StatusError{
		Code:    16,
		Message: "Unauthenticated",
		Details: []proto.Message{
			&ErrorItem{Reason: reason, Message: fmt.Sprintf(format, a)},
		},
	}
}
