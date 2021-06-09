package logging

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func TestHTTP(t *testing.T) {
	var err = errors.New("reply.error")
	var bf = bytes.NewBuffer(nil)
	var logger = log.NewStdLogger(bf)

	tests := []struct {
		name string
		kind func(logger log.Logger) middleware.Middleware
		err  error
		ctx  context.Context
	}{
		{"http-server@fail",
			Server,
			err,
			func() context.Context {
				ctx := transport.NewContext(context.Background(), transport.Transport{Kind: "http", Endpoint: "endpoint"})
				return middleware.WithMethod(ctx, "/package.service/method")
			}(),
		},
		{"http-server@succ",
			Server,
			nil,
			func() context.Context {
				ctx := transport.NewContext(context.Background(), transport.Transport{Kind: "http", Endpoint: "endpoint"})
				return middleware.WithMethod(ctx, "/package.service/method")
			}(),
		},
		{"http-client@succ",
			Client,
			nil,
			func() context.Context {
				ctx := transport.NewContext(context.Background(), transport.Transport{Kind: "http", Endpoint: "endpoint"})
				return middleware.WithMethod(ctx, "/package.service/method")
			}(),
		},
		{"http-client@fail",
			Client,
			err,
			func() context.Context {
				ctx := transport.NewContext(context.Background(), transport.Transport{Kind: "http", Endpoint: "endpoint"})
				return middleware.WithMethod(ctx, "/package.service/method")
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bf.Reset()
			next := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "reply", test.err
			}
			next = test.kind(logger)(next)
			v, e := next(test.ctx, "req.args")
			t.Logf("[%s]reply: %v, error: %v", test.name, v, e)
			t.Logf("[%s]log:%s", test.name, bf.String())
		})
	}
}
