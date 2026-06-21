package errors

import (
	stderrors "errors"
)

// ErrUnsupported indicates that a requested operation cannot be performed,
// because it is unsupported. For example, a call to os.Link when using a
// file system that does not support hard links.
//
// Functions and methods should not return this error but should instead
// return an error including appropriate context that satisfies
//
//	errors.Is(err, errors.ErrUnsupported)
//
// either by directly wrapping ErrUnsupported or by implementing an Is method.
//
// Functions and methods should document the cases in which an error
// wrapping this will be returned.
var ErrUnsupported = stderrors.ErrUnsupported

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target any) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Join returns an error that wraps the given errors. The returned error has a method Unwrap() []error that returns the
// given errors in order.
//
// Join returns nil if errs contains no non-nil error values. If errs contains
// a single non-nil error value, Join returns that error. If errs contains multiple non-nil error values, Join returns an error that formats as the concatenation of the
// format of the non-nil error values, separated by "; ". The returned error's Unwrap method returns a slice of the non-nil error values.
//
// Join is designed for use in situations where multiple errors may be returned, such as when processing a list of items and collecting errors from each item. It allows you to combine those errors into a single error value that can be returned to the caller.
func Join(errs ...error) error {
	return stderrors.Join(errs...)
}
