package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/encoding"
	"github.com/go-kratos/kratos/v3/internal/testdata/binding"
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
