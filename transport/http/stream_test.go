package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/encoding"
	"github.com/go-kratos/kratos/v3/internal/testdata/binding"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/selector"
	transportpkg "github.com/go-kratos/kratos/v3/transport"
)

type streamTestCodec struct{}

func (streamTestCodec) Marshal(any) ([]byte, error) {
	return []byte("stream-test-codec"), nil
}

func (streamTestCodec) Unmarshal(data []byte, v any) error {
	if out, ok := v.(*binding.HelloRequest); ok {
		out.Name = string(data)
	}
	return nil
}

func (streamTestCodec) Name() string {
	return "x-stream-test"
}

func TestServerSentEventStream(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	srv := NewServer()
	ctx := &wrapper{router: srv.Route("/")}
	ctx.Reset(w, req)

	stream := NewServerSentEventServerStream(ctx)
	if err := stream.Send(&binding.HelloRequest{Name: "kratos"}); err != nil {
		t.Fatal(err)
	}
	if err := stream.Close(nil); err != nil {
		t.Fatal(err)
	}

	res := w.Result()
	defer res.Body.Close()
	if got := res.Header.Get("Content-Type"); got != sseContentType {
		t.Fatalf("expected %v, got %v", sseContentType, got)
	}

	clientStream := newSSEClientStream(context.Background(), res, nil)
	var out binding.HelloRequest
	if err := clientStream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "kratos" {
		t.Fatalf("expected %v, got %v", "kratos", out.GetName())
	}
	if err := clientStream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestServerSentEventStreamUsesAcceptCodec(t *testing.T) {
	encoding.RegisterCodec(streamTestCodec{})
	srv := NewServer()
	srv.Route("/").GET("/events", func(ctx Context) error {
		stream := NewServerSentEventServerStream(ctx)
		if err := stream.Send(&binding.HelloRequest{Name: "ignored"}); err != nil {
			return err
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(context.Background(), WithEndpoint(ts.URL), WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.ServerSentEvent(
		context.Background(),
		http.MethodGet,
		"/events",
		nil,
		Accept("text/event-stream, application/x-stream-test"),
		ContentType("application/x-stream-test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	var out binding.HelloRequest
	if err := stream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "stream-test-codec" {
		t.Fatalf("expected custom codec, got %q", out.GetName())
	}
}

func TestServerSentEventStreamUsesClientMiddleware(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/events", func(ctx Context) error {
		if got := ctx.Request().Header.Get("X-Stream-Middleware"); got != "sse" {
			return fmt.Errorf("expected middleware header, got %q", got)
		}
		stream := NewServerSentEventServerStream(ctx)
		if err := stream.Send(&binding.HelloRequest{Name: "kratos"}); err != nil {
			return err
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(
		context.Background(),
		WithEndpoint(ts.URL),
		WithTimeout(time.Second),
		WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				tr, ok := transportpkg.FromClientContext(ctx)
				if !ok {
					return nil, errors.New("missing client transport")
				}
				tr.RequestHeader().Set("X-Stream-Middleware", "sse")
				return handler(ctx, req)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := client.ServerSentEvent(context.Background(), http.MethodGet, "/events", nil, Accept("text/event-stream"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = stream.CloseSend() }()
	var out binding.HelloRequest
	if err := stream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "kratos" {
		t.Fatalf("expected %v, got %v", "kratos", out.GetName())
	}
}

func TestSSEClientStreamClosesBodyOnEOF(t *testing.T) {
	body := &closeCountingBody{}
	res := &http.Response{Header: make(http.Header), Body: body}
	stream := newSSEClientStream(context.Background(), res, nil)

	var out binding.HelloRequest
	if err := stream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
	if body.closed != 1 {
		t.Fatalf("expected body to be closed once, got %d", body.closed)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
	if body.closed != 1 {
		t.Fatalf("expected idempotent close, got %d", body.closed)
	}
}

func TestWebSocketStreamBindsPathQueryAndExchangesMessages(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws/{name}", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}

		in := new(binding.HelloRequest)
		if err := stream.Recv(in); err != nil {
			return stream.Close(err)
		}
		if in.GetName() != "kratos" {
			return stream.Close(fmt.Errorf("expected path name kratos, got %s", in.GetName()))
		}
		if in.GetSub().GetName() != "go" {
			return stream.Close(fmt.Errorf("expected query sub go, got %s", in.GetSub().GetName()))
		}
		if err := stream.Send(&binding.HelloRequest{Name: in.GetName(), Sub: in.GetSub()}); err != nil {
			return stream.Close(err)
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(context.Background(), WithEndpoint(ts.URL), WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	stream, err := client.WebSocket(context.Background(), "/ws/kratos?sub.naming=go", Accept("application/protojson"))
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&binding.HelloRequest{}); err != nil {
		t.Fatal(err)
	}

	var out binding.HelloRequest
	if err := stream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "kratos" {
		t.Fatalf("expected %v, got %v", "kratos", out.GetName())
	}
	if out.GetSub().GetName() != "go" {
		t.Fatalf("expected %v, got %v", "go", out.GetSub().GetName())
	}
	if err := stream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestWebSocketStreamUsesContentTypeCodec(t *testing.T) {
	encoding.RegisterCodec(streamTestCodec{})
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		in := new(binding.HelloRequest)
		if err := stream.Recv(in); err != nil {
			return stream.Close(err)
		}
		if in.GetName() != "stream-test-codec" {
			return stream.Close(fmt.Errorf("expected custom codec, got %q", in.GetName()))
		}
		if err := stream.Send(&binding.HelloRequest{Name: "ignored"}); err != nil {
			return stream.Close(err)
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(context.Background(), WithEndpoint(ts.URL), WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	stream, err := client.WebSocket(
		context.Background(),
		"/ws",
		Accept("application/x-stream-test"),
		ContentType("application/x-stream-test"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&binding.HelloRequest{Name: "ignored"}); err != nil {
		t.Fatal(err)
	}
	var out binding.HelloRequest
	if err := stream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "stream-test-codec" {
		t.Fatalf("expected custom codec, got %q", out.GetName())
	}
}

func TestWebSocketStreamUsesClientMiddleware(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		if got := ctx.Request().Header.Get("X-Stream-Middleware"); got != "websocket" {
			return fmt.Errorf("expected middleware header, got %q", got)
		}
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(
		context.Background(),
		WithEndpoint(ts.URL),
		WithTimeout(time.Second),
		WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (any, error) {
				tr, ok := transportpkg.FromClientContext(ctx)
				if !ok {
					return nil, errors.New("missing client transport")
				}
				tr.RequestHeader().Set("X-Stream-Middleware", "websocket")
				return handler(ctx, req)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	stream, err := client.WebSocket(context.Background(), "/ws")
	if err != nil {
		t.Fatal(err)
	}
	var out binding.HelloRequest
	if err := stream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestWebSocketStreamNormalEOFReportsSelectorSuccess(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		if err := stream.Send(&binding.HelloRequest{Name: "kratos"}); err != nil {
			return stream.Close(err)
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan selector.DoneInfo, 1)
	client, err := NewClient(context.Background(), WithEndpoint(ts.URL), WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	client.r = &resolver{}
	client.selector = &streamTestSelector{
		node: selector.NewNode("http", u.Host, nil),
		done: done,
	}

	stream, err := client.WebSocket(context.Background(), "/ws", Accept("application/protojson"))
	if err != nil {
		t.Fatal(err)
	}
	var out binding.HelloRequest
	if err := stream.Recv(&out); err != nil {
		t.Fatal(err)
	}
	if out.GetName() != "kratos" {
		t.Fatalf("expected %v, got %v", "kratos", out.GetName())
	}
	if err := stream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
	select {
	case di := <-done:
		if di.Err != nil {
			t.Fatalf("expected selector success, got %v", di.Err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for selector done")
	}
}

type closeCountingBody struct {
	closed int
}

func (*closeCountingBody) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (b *closeCountingBody) Close() error {
	b.closed++
	return nil
}

type streamTestSelector struct {
	node selector.Node
	done chan selector.DoneInfo
}

func (s *streamTestSelector) Select(context.Context, ...selector.SelectOption) (selector.Node, selector.DoneFunc, error) {
	return s.node, func(_ context.Context, di selector.DoneInfo) {
		s.done <- di
	}, nil
}

func (*streamTestSelector) Apply([]selector.Node) {}
