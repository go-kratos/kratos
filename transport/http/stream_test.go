package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

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

func TestWebSocketStreamBindsNamedBodyField(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws/{name}", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx, WithStreamBodyField("sub"))
		if err != nil {
			return err
		}

		in := new(binding.HelloRequest)
		if err := stream.Recv(in); err != nil {
			return stream.Close(err)
		}
		// name comes from the path var, the sub message from the streamed frame payload.
		if in.GetName() != "kratos" {
			return stream.Close(fmt.Errorf("expected path name kratos, got %s", in.GetName()))
		}
		if in.GetSub().GetName() != "go" {
			return stream.Close(fmt.Errorf("expected body sub go, got %s", in.GetSub().GetName()))
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

	stream, err := client.WebSocket(context.Background(), "/ws/kratos", Accept("application/protojson"))
	if err != nil {
		t.Fatal(err)
	}
	// The client streams only the body field (Sub), mirroring generated code that
	// sends m.Sub instead of the whole request message.
	if err := stream.Send(&binding.Sub{Name: "go"}); err != nil {
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

func TestServerStreamRecvMessageRejectsInvalidBodyField(t *testing.T) {
	tests := []struct {
		name      string
		bodyField string
	}{
		{name: "unknown field", bodyField: "does_not_exist"},
		{name: "scalar field", bodyField: "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serverStream{mode: streamModeWebSocket, bodyField: tt.bodyField}
			// The field validation happens before any frame is read, so no live
			// connection is required to exercise the error path.
			err := s.recvMessage(new(binding.HelloRequest))
			if err == nil {
				t.Fatalf("expected error for body field %q, got nil", tt.bodyField)
			}
			if !strings.Contains(err.Error(), tt.bodyField) {
				t.Fatalf("expected error to mention %q, got %v", tt.bodyField, err)
			}
		})
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

func TestWebSocketStreamCloseSendPreventsSend(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		in := new(binding.HelloRequest)
		if err := stream.Recv(in); !errors.Is(err, io.EOF) {
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
	stream, err := client.WebSocket(context.Background(), "/ws", Accept("application/protojson"))
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&binding.HelloRequest{Name: "late"}); err == nil {
		t.Fatal("expected Send after CloseSend to fail")
	}

	var out binding.HelloRequest
	if err := stream.Recv(&out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
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

func TestWebSocketStreamSetReadDeadline(t *testing.T) {
	recvErr := make(chan error, 1)
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		if err = stream.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
			recvErr <- err
			return stream.Close(err)
		}
		var in binding.HelloRequest
		err = stream.Recv(&in)
		recvErr <- err
		return stream.Close(err)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	client, err := NewClient(context.Background(), WithEndpoint(ts.URL), WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	// Open the stream but never send, so the server-side Recv hits its read deadline.
	if _, err = client.WebSocket(context.Background(), "/ws", Accept("application/protojson")); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-recvErr:
		if err == nil {
			t.Fatal("expected read deadline error, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for read deadline")
	}
}

func TestWebSocketStreamSetWriteDeadline(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx)
		if err != nil {
			return err
		}
		if err := stream.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
			return stream.Close(err)
		}
		if err := stream.Send(&binding.HelloRequest{Name: "kratos"}); err != nil {
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
}

func TestServerSentEventStreamSetDeadlines(t *testing.T) {
	srv := NewServer()
	var setErr error
	srv.Route("/").GET("/events", func(ctx Context) error {
		stream := NewServerSentEventServerStream(ctx)
		if err := stream.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			setErr = err
		}
		if err := stream.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
			setErr = err
		}
		if err := stream.Send(&binding.HelloRequest{Name: "kratos"}); err != nil {
			return stream.Close(err)
		}
		return stream.Close(nil)
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()
	res, err := ts.Client().Get(ts.URL + "/events")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if _, err := io.ReadAll(res.Body); err != nil {
		t.Fatal(err)
	}
	if setErr != nil {
		t.Fatalf("expected SSE deadline setters to succeed, got %v", setErr)
	}
}

func TestServerStreamSetDeadlineUnknownMode(t *testing.T) {
	s := &serverStream{mode: streamModeWebSocket}
	if err := s.SetReadDeadline(time.Now()); err == nil {
		t.Fatal("expected error when websocket connection is not established")
	}
	if err := s.SetWriteDeadline(time.Now()); err == nil {
		t.Fatal("expected error when websocket connection is not established")
	}
}

type streamCtxKey struct{}

func TestServerStreamDetachesServerTimeout(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	// Simulate the per-request server timeout context plus a middleware-injected value.
	base := context.WithValue(req.Context(), streamCtxKey{}, "trace-123")
	ctx, cancel := context.WithTimeout(base, 10*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	srv := NewServer()
	wrap := &wrapper{router: srv.Route("/")}
	wrap.Reset(w, req)

	stream := NewServerSentEventServerStream(wrap)

	// The server timeout must not propagate to the stream context.
	if _, ok := stream.Context().Deadline(); ok {
		t.Fatal("expected stream context to have no deadline")
	}
	// Middleware-injected values must still be reachable.
	if got := stream.Context().Value(streamCtxKey{}); got != "trace-123" {
		t.Fatalf("expected stream context to preserve values, got %v", got)
	}
	// Even after the per-request timeout fires, the stream context stays alive.
	time.Sleep(20 * time.Millisecond)
	if err := stream.Context().Err(); err != nil {
		t.Fatalf("expected stream context to stay alive after server timeout, got %v", err)
	}

	// SetContext applies the same detachment.
	timedOut, cancel2 := context.WithTimeout(base, time.Nanosecond)
	defer cancel2()
	stream.SetContext(timedOut)
	if _, ok := stream.Context().Deadline(); ok {
		t.Fatal("expected SetContext to detach the deadline")
	}
	if err := stream.Context().Err(); err != nil {
		t.Fatalf("expected stream context to stay alive after SetContext, got %v", err)
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

func TestWebSocketStreamPingKeepalive(t *testing.T) {
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx, WithStreamPing(50*time.Millisecond, 200*time.Millisecond))
		if err != nil {
			return err
		}
		defer func() { _ = stream.Close(nil) }()
		for {
			in := new(binding.HelloRequest)
			if err := stream.Recv(in); err != nil {
				return nil
			}
			if err := stream.Send(in); err != nil {
				return nil
			}
		}
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()

	conn, res, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	// On a successful upgrade gorilla replaces the response body with a
	// no-op reader; closing it is safe and keeps bodyclose happy.
	defer func() { _ = res.Body.Close() }()
	defer conn.Close()

	pingSeen := make(chan struct{})
	var pingOnce sync.Once
	conn.SetPingHandler(func(appData string) error {
		pingOnce.Do(func() { close(pingSeen) })
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})
	readErr := make(chan error, 1)
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				readErr <- err
				return
			}
			readErr <- errors.New(string(data))
		}
	}()

	select {
	case <-pingSeen:
	case <-time.After(2 * time.Second):
		t.Fatal("expected a server ping")
	}

	// Stay silent well beyond pongWait: the connection must survive because the
	// client's pongs keep extending the server-side read deadline.
	time.Sleep(500 * time.Millisecond)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"name":"alive"}`)); err != nil {
		t.Fatalf("write after pong keepalive: %v", err)
	}
	select {
	case err := <-readErr:
		if err == nil || !strings.Contains(err.Error(), "alive") {
			t.Fatalf("expected echoed message, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected echo after keepalive")
	}
}

func TestWebSocketStreamPingDropsSilentPeer(t *testing.T) {
	recvErr := make(chan error, 1)
	srv := NewServer()
	srv.Route("/").GET("/ws", func(ctx Context) error {
		stream, err := NewWebSocketServerStream(ctx, WithStreamPing(50*time.Millisecond, 120*time.Millisecond))
		if err != nil {
			return err
		}
		defer func() { _ = stream.Close(nil) }()
		for {
			if err := stream.Recv(new(binding.HelloRequest)); err != nil {
				recvErr <- err
				return nil
			}
		}
	})

	ts := httptest.NewServer(srv)
	defer ts.Close()

	conn, res, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(ts.URL, "http")+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	// On a successful upgrade gorilla replaces the response body with a
	// no-op reader; closing it is safe and keeps bodyclose happy.
	defer func() { _ = res.Body.Close() }()
	defer conn.Close()
	// Swallow pings without ever ponging, simulating a half-open peer.
	conn.SetPingHandler(func(string) error { return nil })
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	select {
	case err := <-recvErr:
		if err == nil {
			t.Fatal("expected server Recv to fail for a peer that never pongs")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected server Recv to fail for a peer that never pongs")
	}
}
