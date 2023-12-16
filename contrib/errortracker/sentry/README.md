# Sentry middleware for Kratos
This middleware helps you to catch panics and report them to [sentry](https://sentry.io/)

## Quick Start
You could check the full demo in example folder.
```go
// Step 1: 
// init sentry in the entry of your application
import "github.com/getsentry/sentry-go"

sentry.Init(sentry.ClientOptions{
		Dsn: "<your dsn>",
		AttachStacktrace: true, // recommended
})


// Step 2: 
// set middleware
import 	ksentry "github.com/go-kratos/sentry"

// for HTTP server, new HTTP server with sentry middleware options
var opts = []http.ServerOption{
	http.Middleware(
		recovery.Recovery(), 
		ksentry.Server(ksentry.WithTags(map[string]interface{}{
			"tag": "some-custom-constant-tag",
			"trace_id": tracing.TraceID(),
		})), // must after Recovery middleware, because of the exiting order will be reversed
		tracing.Server(),
		logging.Server(logger), 
	),
}

// for gRPC server, new gRPC server with sentry middleware options
var opts = []grpc.ServerOption{
     grpc.Middleware(
		recovery.Recovery(),
		ksentry.Server(ksentry.WithTags(map[string]interface{}{
			"tag": "some-custom-constant-tag",
			"trace_id": tracing.TraceID(),
		})), // must after Recovery middleware, because of the exiting order will be reversed
		tracing.Server(),
		logging.Server(logger),
     ),
 }


// Then, the framework will report events to Sentry when your trigger panics.
// Or your can push events to Sentry manually
```

## Reference
* [https://docs.sentry.io/platforms/go/](https://docs.sentry.io/platforms/go/)