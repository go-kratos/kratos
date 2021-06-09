package metadata

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	gmetadata "google.golang.org/grpc/metadata"
)

// Server is an server metadata middleware.
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			md := metadata.New()

			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					info, _ := http.FromServerContext(ctx)
					if info.Request != nil {
						for key, values := range info.Request.Header {
							key = strings.ToLower(key)
							if len(values) == 1 && strings.HasPrefix(key, "x-md-") {
								md.Set(strings.TrimLeft(key, "x-md-"), values[0])
							}
						}
					}
				case transport.KindGRPC:
					gmd, _ := gmetadata.FromIncomingContext(ctx)
					for key, values := range gmd {
						if len(values) == 1 {
							md.Set(key, values[0])
						}
					}
				}
			}
			ctx = metadata.NewContext(ctx, md)

			reply, err = handler(ctx, req)
			return
		}
	}
}

// Client is an server metadata middleware.
func Client() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			md, _ := metadata.FromContext(ctx)

			if tr, ok := transport.FromContext(ctx); ok {
				switch tr.Kind {
				case transport.KindHTTP:
					info, _ := http.FromClientContext(ctx)
					if info.Request != nil {
						md.Range(func(k, v string) bool {
							info.Request.Header.Set("x-md-"+k, v)
							return true
						})
					}
				case transport.KindGRPC:
					gmd, _ := gmetadata.FromOutgoingContext(ctx)
					// copy md to avoid datarace
					gmd = gmd.Copy()
					md.Range(func(k, v string) bool {
						gmd.Set(k, v)
						return true
					})
					if len(gmd) > 0 {
						ctx = gmetadata.NewOutgoingContext(ctx, gmd)
					}
				}
			}

			reply, err = handler(ctx, req)
			return
		}
	}
}
