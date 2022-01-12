package jaeger

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/net/trace"
)

var (
	_jaegerAppID    = env.AppID
	_jaegerEndpoint = "http://127.0.0.1:9191"
	_probability    = 0.00025
)

func init() {
	if v := os.Getenv("JAEGER_ENDPOINT"); v != "" {
		_jaegerEndpoint = v
	}

	if v := os.Getenv("JAEGER_APPID"); v != "" {
		_jaegerAppID = v
	}

	flag.StringVar(&_jaegerEndpoint, "jaeger_endpoint", _jaegerEndpoint, "jaeger report endpoint, or use JAEGER_ENDPOINT env.")
	flag.StringVar(&_jaegerAppID, "jaeger_appid", _jaegerAppID, "jaeger report appid, or use JAEGER_APPID env.")
}

// Init Init
func Init(cfg *Config) {
	c := cfg
	if c == nil {
		c = &Config{AppID: _jaegerAppID, Endpoint: _jaegerEndpoint, BatchSize: 120, Probability: float32(_probability)}
	}
	trace.SetGlobalTracer(trace.NewTracer(c.AppID, newReport(c), true, c.Probability))
}
