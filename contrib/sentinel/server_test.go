package sentinel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
)

func initServerSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "GET:/ping",
			Threshold:              1.0,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "/api/123",
			Threshold:              0.0,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
		return
	}
}

func TestSentinelServerMiddleware(t *testing.T) {
	type args struct {
		opts    []Option
		method  string
		path    string
		reqPath string
		handler khttp.HandlerFunc
		body    io.Reader
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "default get",
			args: args{
				opts:    []Option{},
				method:  http.MethodGet,
				path:    "/ping",
				reqPath: "/ping",
				handler: khttp.HandlerFunc(func(ctx khttp.Context) error {
					khttp.SetOperation(ctx, "/test/ping")
					h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
						return "ping", nil
					})
					// may use `ctx.BindQuery()` and `ctx.BindVars()` to build `req` for `h`
					out, err := h(ctx, nil)
					if err != nil {
						return err
					}
					return ctx.Result(200, fmt.Sprintf("%+v", out))
				}),
				body: nil,
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "customize resource extract",
			args: args{
				opts: []Option{
					WithResourceExtractor(func(ctx context.Context, req interface{}) string {
						tr, ok := transport.FromServerContext(ctx)
						if ok {
							httpTr := tr.(khttp.Transporter)
							return httpTr.Request().URL.Path
						}
						return ""
					}),
				},
				method:  http.MethodGet,
				path:    "/api/{uid}",
				reqPath: "/api/123",
				handler: khttp.HandlerFunc(func(ctx khttp.Context) error {
					khttp.SetOperation(ctx, "/test/api")
					h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
						return "api", nil
					})
					out, err := h(ctx, nil)
					if err != nil {
						return err
					}
					return ctx.Result(200, fmt.Sprintf("%+v", out))
				}),
				body: nil,
			},
			want: want{
				code: http.StatusTooManyRequests,
			},
		},
		{
			name: "customize block fallback",
			args: args{
				opts: []Option{
					WithBlockFallback(func(ctx context.Context, req interface{}) (interface{}, error) {
						return nil, errors.New(http.StatusBadRequest, "Customized Error", "Blocked by Sentinel")
					}),
				},
				method:  http.MethodGet,
				path:    "/ping",
				reqPath: "/ping",
				handler: khttp.HandlerFunc(func(ctx khttp.Context) error {
					khttp.SetOperation(ctx, "/test/ping")
					h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
						return "ping", nil
					})
					out, err := h(ctx, nil)
					if err != nil {
						return err
					}
					return ctx.Result(200, fmt.Sprintf("%+v", out))
				}),
				body: nil,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	initServerSentinel(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []khttp.ServerOption{
				khttp.Middleware(
					ServerMiddleware(tt.args.opts...),
				),
			}
			httpSrv := khttp.NewServer(opts...)
			router := httpSrv.Route("/")
			router.GET(tt.args.path, tt.args.handler)

			r := httptest.NewRequest(tt.args.method, tt.args.reqPath, nil)
			w := httptest.NewRecorder()
			httpSrv.ServeHTTP(w, r)
			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
