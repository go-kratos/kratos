package sentinel

import (
	"context"
	"fmt"
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// ServerMiddleware returns new middleware.Handler for kratos http/grpc server.
// Default resource name pattern is {httpMethod}:{apiPath}, such as "GET:/api/:id".
// Default block fallback is to return 429 (Too Many Requests) response.
//
// You may customize your own resource extractor and block handler by setting options.
func ServerMiddleware(opts ...Option) middleware.Middleware {
	options := evaluateOptions(opts)
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var resourceName string
			var resType base.ResourceType
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				err = errors.New(http.StatusBadRequest, "Bad Request", "failed to extract request from context")
				return nil, err
			}
			if options.resourceExtract != nil {
				resourceName = options.resourceExtract(ctx, req)
			} else {
				switch tr.Kind() {
				case transport.KindGRPC:
					resourceName = tr.Operation()
					resType = base.ResTypeRPC
				case transport.KindHTTP:
					httpTr := tr.(khttp.Transporter)
					resourceName = fmt.Sprintf("%s:%s", httpTr.Request().Method, httpTr.PathTemplate())
					resType = base.ResTypeWeb
				default:
					err = errors.New(http.StatusBadRequest, "Bad Request", fmt.Sprintf("unsupported transport kind: %s", tr.Kind()))
					return nil, err
				}
			}
			// start building sentinel entry
			entry, blockErr := sentinel.Entry(
				resourceName,
				sentinel.WithResourceType(resType),
				sentinel.WithTrafficType(base.Inbound),
			)
			if blockErr != nil {
				if options.blockFallback != nil {
					reply, err = options.blockFallback(ctx, req)
					return
				}
				switch tr.Kind() {
				case transport.KindGRPC:
					err = blockErr
				case transport.KindHTTP:
					err = errors.New(http.StatusTooManyRequests, "Too many requests", blockErr.Error())
				}
				return nil, err
			}
			defer entry.Exit()

			reply, err = next(ctx, req)
			return
		}
	}
}
