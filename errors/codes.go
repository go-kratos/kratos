package errors

import (
	"errors"
	"fmt"
)

// Cancelled The operation was cancelled, typically by the caller.
// HTTP Mapping: 499 Client Closed Request
func Cancelled(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    1,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsCancelled(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 1
	}
	return false
}

// Unknown error.
// HTTP Mapping: 500 Internal Server Error
func Unknown(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    2,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsUnknown(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 2
	}
	return false
}

// InvalidArgument The client specified an invalid argument.
// HTTP Mapping: 400 Bad Request
func InvalidArgument(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    3,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsInvalidArgument(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 3
	}
	return false
}

// DeadlineExceeded The deadline expired before the operation could complete.
// HTTP Mapping: 504 Gateway Timeout
func DeadlineExceeded(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    4,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsDeadlineExceeded(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 4
	}
	return false
}

// NotFound Some requested entity (e.g., file or directory) was not found.
// HTTP Mapping: 404 Not Found
func NotFound(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    5,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsNotFound(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 5
	}
	return false
}

// AlreadyExists The entity that a client attempted to create (e.g., file or directory) already exists.
// HTTP Mapping: 409 Conflict
func AlreadyExists(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    6,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsAlreadyExists(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 6
	}
	return false
}

// PermissionDenied The caller does not have permission to execute the specified operation.
// HTTP Mapping: 403 Forbidden
func PermissionDenied(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    7,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsPermissionDenied(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 7
	}
	return false
}

// ResourceExhausted Some resource has been exhausted, perhaps a per-user quota, or
// perhaps the entire file system is out of space.
// HTTP Mapping: 429 Too Many Requests
func ResourceExhausted(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    8,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsResourceExhausted(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 8
	}
	return false
}

// FailedPrecondition The operation was rejected because the system is not in a state
// required for the operation's execution.
// HTTP Mapping: 400 Bad Request
func FailedPrecondition(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    9,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsFailedPrecondition(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 9
	}
	return false
}

// Aborted The operation was aborted, typically due to a concurrency issue such as
// a sequencer check failure or transaction abort.
// HTTP Mapping: 409 Conflict
func Aborted(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    10,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsAborted(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 10
	}
	return false
}

// OutOfRange The operation was attempted past the valid range.  E.g., seeking or
// reading past end-of-file.
// HTTP Mapping: 400 Bad Request
func OutOfRange(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    11,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsOutOfRange(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 11
	}
	return false
}

// Unimplemented The operation is not implemented or is not supported/enabled in this service.
// HTTP Mapping: 501 Not Implemented
func Unimplemented(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    12,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsUnimplemented(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 12
	}
	return false
}

// Internal This means that some invariants expected by the
// underlying system have been broken.  This error code is reserved
// for serious errors.
//
// HTTP Mapping: 500 Internal Server Error
func Internal(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    13,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsInternal(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 13
	}
	return false
}

// Unavailable The service is currently unavailable.
// HTTP Mapping: 503 Service Unavailable
func Unavailable(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    14,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsUnavailable(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 14
	}
	return false
}

// DataLoss Unrecoverable data loss or corruption.
// HTTP Mapping: 500 Internal Server Error
func DataLoss(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    15,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsDataLoss(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 15
	}
	return false
}

// Unauthorized The request does not have valid authentication credentials for the operation.
// HTTP Mapping: 401 Unauthorized
func Unauthorized(reason, format string, a ...interface{}) error {
	return &StatusError{
		Code:    16,
		Reason:  reason,
		Message: fmt.Sprintf(format, a...),
	}
}

func IsUnauthorized(err error) bool {
	if se := new(StatusError); errors.As(err, &se) {
		return se.Code == 16
	}
	return false
}
