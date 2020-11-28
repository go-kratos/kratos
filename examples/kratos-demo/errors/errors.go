package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
	apierr "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/errors"
)

// MissingField is A required field on a resource has not been set.
func MissingField() error {
	return errors.InvalidArgument(
		apierr.Kratos_MissingField.String(),
		"name is missing",
	)
}

// RequestBlocked is the requesting user was blocked.
func RequestBlocked(message string) error {
	return errors.InvalidArgument(
		apierr.Kratos_RequestBlocked.String(),
		"request is blocked",
	)
}
