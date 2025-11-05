package http

import (
	"context"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
)

func TestTransport_Kind(t *testing.T) {
	o := &Transport{}
	if !reflect.DeepEqual(transport.KindHTTP, o.Kind()) {
		t.Errorf("expect %v, got %v", transport.KindHTTP, o.Kind())
	}
}

func TestTransport_Endpoint(t *testing.T) {
	v := "hello"
	o := &Transport{endpoint: v}
	if !reflect.DeepEqual(v, o.Endpoint()) {
		t.Errorf("expect %v, got %v", v, o.Endpoint())
	}
}

func TestTransport_Operation(t *testing.T) {
	v := "hello"
	o := &Transport{operation: v}
	if !reflect.DeepEqual(v, o.Operation()) {
		t.Errorf("expect %v, got %v", v, o.Operation())
	}
}

func TestTransport_Request(t *testing.T) {
	v := &http.Request{}
	o := &Transport{request: v}
	if !reflect.DeepEqual(v, o.Request()) {
		t.Errorf("expect %v, got %v", v, o.Request())
	}
}

func TestTransport_RequestHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{reqHeader: v}
	if !reflect.DeepEqual("1", o.RequestHeader().Get("a")) {
		t.Errorf("expect %v, got %v", "1", o.RequestHeader().Get("a"))
	}
}

func TestTransport_Response(t *testing.T) {
	v := http.ResponseWriter(nil)
	o := &Transport{response: v}
	if !reflect.DeepEqual(v, o.Response()) {
		t.Errorf("expect %v, got %v", v, o.Response())
	}
}

func TestTransport_ReplyHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{replyHeader: v}
	if !reflect.DeepEqual("1", o.ReplyHeader().Get("a")) {
		t.Errorf("expect %v, got %v", "1", o.ReplyHeader().Get("a"))
	}
}

func TestTransport_PathTemplate(t *testing.T) {
	v := "template"
	o := &Transport{pathTemplate: v}
	if !reflect.DeepEqual(v, o.PathTemplate()) {
		t.Errorf("expect %v, got %v", v, o.PathTemplate())
	}
}

func TestHeaderCarrier_Keys(t *testing.T) {
	v := headerCarrier{}
	v.Set("abb", "1")
	v.Set("bcc", "2")
	want := []string{"Abb", "Bcc"}
	keys := v.Keys()
	sort.Slice(want, func(i, j int) bool {
		return want[i] < want[j]
	})
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	if !reflect.DeepEqual(want, keys) {
		t.Errorf("expect %v, got %v", want, keys)
	}
}

func TestSetOperation(t *testing.T) {
	tr := &Transport{}
	ctx := transport.NewServerContext(context.Background(), tr)
	SetOperation(ctx, "kratos")
	if !reflect.DeepEqual(tr.operation, "kratos") {
		t.Errorf("expect %v, got %v", "kratos", tr.operation)
	}
}

// TestResponseTransporter_Interface tests that Transport implements ResponseTransporter
func TestResponseTransporter_Interface(t *testing.T) {
	var transport Transporter = &Transport{}
	if _, ok := transport.(ResponseTransporter); !ok {
		t.Error("Transport should implement ResponseTransporter interface")
	}
}

// TestResponseWriterFromServerContext tests the ResponseWriterFromServerContext helper function
func TestResponseWriterFromServerContext(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func() context.Context
		expectWriter bool
		expectOk     bool
	}{
		{
			name: "valid HTTP transport with ResponseWriter",
			setupContext: func() context.Context {
				mockWriter := &mockResponseWriter{header: make(http.Header)}
				tr := &Transport{response: mockWriter}
				return transport.NewServerContext(context.Background(), tr)
			},
			expectWriter: true,
			expectOk:     true,
		},
		{
			name: "context without transport",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectWriter: false,
			expectOk:     false,
		},
		{
			name: "context with non-HTTP transport",
			setupContext: func() context.Context {
				tr := &mockNonHTTPTransport{}
				return transport.NewServerContext(context.Background(), tr)
			},
			expectWriter: false,
			expectOk:     false,
		},
		{
			name: "context with HTTP transport without ResponseTransporter interface",
			setupContext: func() context.Context {
				tr := &mockBasicHTTPTransport{}
				return transport.NewServerContext(context.Background(), tr)
			},
			expectWriter: false,
			expectOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			writer, ok := ResponseWriterFromServerContext(ctx)

			if ok != tt.expectOk {
				t.Errorf("ResponseWriterFromServerContext() ok = %v, want %v", ok, tt.expectOk)
			}

			if tt.expectWriter && writer == nil {
				t.Error("ResponseWriterFromServerContext() should return non-nil writer")
			}

			if !tt.expectWriter && writer != nil {
				t.Error("ResponseWriterFromServerContext() should return nil writer")
			}
		})
	}
}

// TestTransport_InterfaceCompatibility tests interface compatibility
func TestTransport_InterfaceCompatibility(t *testing.T) {
	tr := &Transport{
		endpoint:     "http://localhost:8080",
		operation:    "/test",
		pathTemplate: "/test/{id}",
		request:      &http.Request{},
		response:     &mockResponseWriter{header: make(http.Header)},
	}

	// Test basic Transporter interface
	var basicTransporter Transporter = tr
	if basicTransporter.Request() == nil {
		t.Error("Transporter.Request() should not be nil")
	}
	if basicTransporter.PathTemplate() == "" {
		t.Error("Transporter.PathTemplate() should not be empty")
	}

	// Test ResponseTransporter interface
	var responseTransporter ResponseTransporter = tr
	if responseTransporter.Response() == nil {
		t.Error("ResponseTransporter.Response() should not be nil")
	}

	// Test that ResponseTransporter extends Transporter
	if responseTransporter.Request() == nil {
		t.Error("ResponseTransporter should have Request() method from Transporter")
	}
}

// TestTransport_TypeAssertion tests safe type assertion patterns
func TestTransport_TypeAssertion(t *testing.T) {
	tr := &Transport{
		response: &mockResponseWriter{header: make(http.Header)},
	}
	ctx := transport.NewServerContext(context.Background(), tr)

	// Test type assertion to ResponseTransporter
	if serverTr, ok := transport.FromServerContext(ctx); ok {
		if httpTr, ok := serverTr.(ResponseTransporter); ok {
			if httpTr.Response() == nil {
				t.Error("ResponseTransporter.Response() should not be nil")
			}
		} else {
			t.Error("Transport should be assertable to ResponseTransporter")
		}
	} else {
		t.Error("Should be able to get transport from server context")
	}
}

// Mock implementations for testing
// Note: mockResponseWriter is already defined in codec_test.go and will be reused

// mockNonHTTPTransport simulates a non-HTTP transport (like gRPC)
type mockNonHTTPTransport struct{}

func (m *mockNonHTTPTransport) Kind() transport.Kind            { return transport.KindGRPC }
func (m *mockNonHTTPTransport) Endpoint() string                { return "grpc://localhost:9000" }
func (m *mockNonHTTPTransport) Operation() string               { return "/grpc.Service/Method" }
func (m *mockNonHTTPTransport) RequestHeader() transport.Header { return nil }
func (m *mockNonHTTPTransport) ReplyHeader() transport.Header   { return nil }

// mockBasicHTTPTransport simulates an HTTP transport that only implements basic Transporter
type mockBasicHTTPTransport struct{}

func (m *mockBasicHTTPTransport) Kind() transport.Kind            { return transport.KindHTTP }
func (m *mockBasicHTTPTransport) Endpoint() string                { return "http://localhost:8080" }
func (m *mockBasicHTTPTransport) Operation() string               { return "/test" }
func (m *mockBasicHTTPTransport) RequestHeader() transport.Header { return nil }
func (m *mockBasicHTTPTransport) ReplyHeader() transport.Header   { return nil }
func (m *mockBasicHTTPTransport) Request() *http.Request          { return nil }
func (m *mockBasicHTTPTransport) PathTemplate() string            { return "/test" }
