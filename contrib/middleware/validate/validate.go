package validate

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

type validator interface {
	Validate() error
}

// ProtoValidate is a middleware that validates the request message with [protovalidate](https://github.com/bufbuild/protovalidate)
func ProtoValidate() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if msg, ok := req.(proto.Message); ok {
				if err := protovalidate.Validate(msg); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}

			// to compatible with the [old validator](https://github.com/envoyproxy/protoc-gen-validate)
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}

			return handler(ctx, req)
		}
	}
}
