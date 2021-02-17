package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
	"github.com/go-kratos/kratos/v2/middleware"
)

type testRequest struct {
	Name string `json:"name"`
}
type testReply struct {
	Result string `json:"result"`
}
type testService struct{}

func (s *testService) SayHello(ctx context.Context, req *testRequest) (*testReply, error) {
	return &testReply{Result: req.Name}, nil
}

func TestService(t *testing.T) {
	h := func(srv interface{}, ctx context.Context, req *http.Request, dec func(interface{}) error, m middleware.Middleware) (interface{}, error) {
		var in testRequest
		if err := dec(&in); err != nil {
			return nil, err
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(*testService).SayHello(ctx, &in)
		}
		out, err := m(h)(ctx, &in)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	sd := &ServiceDesc{
		ServiceName: "helloworld.Greeter",
		Methods: []MethodDesc{
			{
				Path:    "/helloworld",
				Method:  "POST",
				Handler: h,
			},
		},
	}

	svc := &testService{}
	srv := NewServer()
	srv.RegisterService(sd, svc)

	time.AfterFunc(time.Second, func() {
		defer srv.Stop()
		testServiceClient(t, srv)
	})

	if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
		t.Fatal(err)
	}
}

func testServiceClient(t *testing.T, srv *Server) {
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	port, ok := host.Port(srv.lis)
	if !ok {
		t.Fatalf("extract port error: %v", srv.lis)
	}
	var (
		in  = testRequest{Name: "hello"}
		out = testReply{}
		url = fmt.Sprintf("http://127.0.0.1:%d/helloworld", port)
	)
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/json")
	if err := Do(client, req, &out); err != nil {
		t.Fatal(err)
	}
	if out.Result != in.Name {
		t.Fatalf("expected %s got %s", in.Name, out.Result)
	}
}
