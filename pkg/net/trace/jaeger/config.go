package jaeger

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/net/trace"
)

var _jaegerEndpoint = "http://127.0.0.1:9191"

func init() {
	if v := os.Getenv("JAEGER_ENDPOINT"); v != "" {
		_jaegerEndpoint = v
	}
	flag.StringVar(&_jaegerEndpoint, "jaeger_endpoint", _jaegerEndpoint, "jaeger report endpoint, or use JAEGER_ENDPOINT env.")
}

// Init Init
func Init() {
	c := &Config{Endpoint: _jaegerEndpoint, BatchSize: 120}
	trace.SetGlobalTracer(trace.NewTracer(env.AppID, newReport(c), true))
}
