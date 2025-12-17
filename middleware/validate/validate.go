package validate

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// ValidateFunc defines a validation function type.
type ValidateFunc func(v any) error

// requestValidator is an interface for types that can validate themselves.
type requestValidator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator(validators ...ValidateFunc) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if v, ok := req.(requestValidator); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			for _, v := range validators {
				if err := v(req); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
