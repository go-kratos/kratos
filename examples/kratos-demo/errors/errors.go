package errors

import (
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/codes"
	"github.com/go-kratos/kratos/v2/status"
)

var (
	domain = "go-kratos.dev"

	// ErrMissingField is A required field on a resource has not been set.
	ErrMissingField = status.BadRequest(
		&status.ErrorInfo{
			Reason: codes.Kratos_MissingField.String(),
			Domain: domain,
		},
	)
)

// RequestBlocked is the requesting user was blocked.
func RequestBlocked(resouces string) error {
	return status.InternalServerError(
		&status.ErrorInfo{
			Reason:   codes.Kratos_RequestBlocked.String(),
			Domain:   domain,
			Metadata: map[string]string{"resouces": resouces},
		},
	)
}
