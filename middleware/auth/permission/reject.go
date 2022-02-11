package permission

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

const (
	// reason holds the error reason.
	reason = "NO_PERMISSION"

	// message holds the error message.
	message = "You do not have permission to access"
)

var (
	// NoPermission holds NoPermission error
	NoPermission = errors.Forbidden(reason, message)
)

// server . This function just reject.
func server() middleware.Middleware {

	// reject Middleware.
	return func(handler middleware.Handler) middleware.Handler {
		// reject Handler
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// just reject and return no permission.
			return nil, NoPermission
		}
	}

}
