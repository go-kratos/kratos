package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/codes"
)

var (
	// ErrMissingField is A required field on a resource has not been set.
	ErrMissingField = errors.InvalidArgument(
		codes.Kratos_MissingField.String(),
		"name is missing",
	)
)

// RequestBlocked is the requesting user was blocked.
func RequestBlocked(message string) error {
	return errors.Internal(
		codes.Kratos_RequestBlocked.String(),
		message,
	)
}
