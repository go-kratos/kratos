/*
This package provides Sentinel integration for Kratos.

Kratos provides unified abstraction for grpc and http middlewares.
Here we take the server side as examples.

For kratos/transport/http, user can call `http.Middleware`, for example,

	import (
		"github.com/go-kratos/kratos/v2/transport/http"
		"github.com/go-kratos/kratos/contrib/sentinel"
	)
	var opts = []http.ServerOption{
		http.Middleware(
			sentinel.ServerMiddleware(),
		),
	}
	server := http.NewServer(opts...)

For kratos/transport/grpc, user can call `grpc.Middleware`, for example,

	import (
		"github.com/go-kratos/kratos/v2/transport/grpc"
		"github.com/go-kratos/kratos/contrib/sentinel"
	)
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			sentinel.ServerMiddleware(),
		),
	}
	server := grpc.NewServer(opts...)

User can also use sentinel grpc interceptors, for example,

	import (
		"github.com/go-kratos/kratos/v2/transport/grpc"
		sentinel "github.com/sentinel-go/pkg/adapters/grpc"
	)
	var opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(
			sentinel.NewUnaryServerInterceptor(),
		),
	}
	server := grpc.NewServer(opts...)

The plugin extracts `Request().Method:PathTemplate()` (for http) or `Operation()`
(for grpc) as the resource name by default. Users may provide customized
resource name extractor when creating new Sentinel middlewares (via options).

Fallback logic: the plugin will return the BlockError by default
if current request is blocked by Sentinel rules. Users may also
provide customized fallback logic via WithXxxBlockFallback(handler) options.
*/
package sentinel
