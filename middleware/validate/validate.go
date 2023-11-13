package validate

import (
	"context"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ugly but permit to allocate one time CLE engine
// https://github.com/bufbuild/protovalidate-go/issues/71#issuecomment-1781624431
var validator *protovalidate.Validator

func init() {
	var err error

	if validator, err = protovalidate.New(); err != nil {
		fmt.Println("failed to initialize validator:", err)
		os.Exit(1)
	}
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if err := validator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
				return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
			}
			return handler(ctx, req)
		}
	}
}
