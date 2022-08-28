package sentinel

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
)

func initClientSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "GET:/client",
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
		{
			Resource:               "GET:/custom_fallback",
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

func TestSentinelClientMiddleware(t *testing.T) {
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
				path:    "/client",
				reqPath: "/client",
				handler: khttp.HandlerFunc(func(ctx khttp.Context) error {
					khttp.SetOperation(ctx, "/test/client")
					h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
						return "client", nil
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
						tr, ok := transport.FromClientContext(ctx)
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
				path:    "/custom_fallback",
				reqPath: "/custom_fallback",
				handler: khttp.HandlerFunc(func(ctx khttp.Context) error {
					khttp.SetOperation(ctx, "/test/custom_fallback")
					h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
						return "custom_fallback", nil
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
	initClientSentinel(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, _ := getAvailablePort(8000)
			// create server
			srvOpts := []khttp.ServerOption{
				khttp.Address(fmt.Sprintf(":%d", port)),
			}
			httpSrv := khttp.NewServer(srvOpts...)
			router := httpSrv.Route("/")
			router.GET(tt.args.path, tt.args.handler)
			go func() { _ = httpSrv.Start(context.Background()) }()
			defer func() { _ = httpSrv.Stop(context.Background()) }()
			time.Sleep(time.Second)
			// create client
			connOpts := []khttp.ClientOption{
				khttp.WithMiddleware(
					ClientMiddleware(tt.args.opts...),
				),
				khttp.WithEndpoint(fmt.Sprintf("127.0.0.1:%d", port)),
			}
			conn, _ := khttp.NewClient(context.Background(), connOpts...)
			defer conn.Close()
			// invoke request
			var out interface{}
			callOpts := []khttp.CallOption{
				khttp.Operation(tt.args.path),
				khttp.PathTemplate(tt.args.path),
			}
			err := conn.Invoke(context.Background(), tt.args.method, tt.args.reqPath, nil, &out, callOpts...)
			assert.Equal(t, tt.want.code, errors.Code(err))
		})
	}
}

func getAvailablePort(init int) (int, error) {
	for p := init; p < 65536; p++ {
		conn, _ := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprint(p)), time.Second)
		if conn != nil {
			conn.Close()
		} else {
			return p, nil
		}
	}
	return 0, fmt.Errorf("Cannot get an available port")
}
