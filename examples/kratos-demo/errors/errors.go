package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
	code "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/errors"
)

var (
	domain = "api.go-kratos.dev"
)

func WithErrorInfo(reason string, md map[string]string) error {
	return &errors.ErrorInfo{
		Reason:   reason,
		Domain:   domain,
		Metadata: md,
	}
}

// ErrMissingField is A required field on a resource has not been set.
func ErrMissingField(md map[string]string) error {
	return errors.BadRequest(
		errors.ErrorInfo{
			Reason:   code.Errors_MissingField.String(),
			Domain:   domain,
			Metadata: md,
		},
	)
}
