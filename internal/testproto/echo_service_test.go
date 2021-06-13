package testproto

import (
	context "context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/metadata"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_struct "github.com/golang/protobuf/ptypes/struct"
)

var md = metadata.Metadata{"x-md-global-test": "test_value"}

type echoService struct {
	UnimplementedEchoServiceServer
}

func (s *echoService) Echo(ctx context.Context, m *SimpleMessage) (*SimpleMessage, error) {
	md, _ := metadata.FromServerContext(ctx)
	if v := md.Get("x-md-global-test"); v != "test_value" {
		return nil, errors.New("md not match" + v)
	}
	return m, nil
}

func (s *echoService) EchoBody(ctx context.Context, m *SimpleMessage) (*SimpleMessage, error) {
	return m, nil
}

func (s *echoService) EchoDelete(ctx context.Context, m *SimpleMessage) (*SimpleMessage, error) {
	return m, nil
}

func (s *echoService) EchoPatch(ctx context.Context, m *DynamicMessageUpdate) (*DynamicMessageUpdate, error) {
	return m, nil
}

func (s *echoService) EchoResponseBody(ctx context.Context, m *DynamicMessageUpdate) (*DynamicMessageUpdate, error) {
	return m, nil
}

type echoClient struct {
	client EchoServiceHTTPClient
}

// post: /v1/example/echo/{id}
func (c *echoClient) Echo(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	return c.client.Echo(ctx, in)
}

// post: /v1/example/echo_body
func (c *echoClient) EchoBody(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	return c.client.EchoBody(ctx, in)
}

// delete: /v1/example/echo_delete/{id}/{num}
func (c *echoClient) EchoDelete(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	return c.client.EchoDelete(ctx, in)
}

// patch: /v1/example/echo_patch
func (c *echoClient) EchoPatch(ctx context.Context, in *DynamicMessageUpdate) (out *DynamicMessageUpdate, err error) {
	return c.client.EchoPatch(ctx, in)
}

// post: /v1/example/echo_response_body
func (c *echoClient) EchoResponseBody(ctx context.Context, in *DynamicMessageUpdate) (out *DynamicMessageUpdate, err error) {
	return c.client.EchoResponseBody(ctx, in)
}

func TestJSON(t *testing.T) {
	in := &SimpleMessage{Id: "test_id", Num: 100}
	out := &SimpleMessage{}
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if err := codec.Unmarshal(data, out); err != nil {
		t.Fatal(err)
	}
	// body
	in2 := &DynamicMessageUpdate{Body: &DynamicMessage{
		ValueField: &_struct.Value{Kind: &_struct.Value_StringValue{StringValue: "test"}},
	}}
	out2 := &DynamicMessageUpdate{}
	data, err = codec.Marshal(&in2.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := codec.Unmarshal(data, &out2.Body); err != nil {
		t.Fatal(err)
	}
}

func TestEchoHTTPServer(t *testing.T) {
	echo := &echoService{}
	ctx := context.Background()
	srv := http.NewServer(
		http.Address(":2333"),
		http.Middleware(mmd.Server()),
	)
	RegisterEchoServiceHTTPServer(srv, echo)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testEchoHTTPClient(t, fmt.Sprintf("127.0.0.1:2333"))
	srv.Stop(ctx)
}

func testEchoHTTPClient(t *testing.T, addr string) {
	var (
		err error
		in  = &SimpleMessage{Id: "test_id", Num: 100}
		out = &SimpleMessage{}
	)
	check := func(name string, in, out *SimpleMessage) {
		if in.Id != out.Id || in.Num != out.Num {
			t.Errorf("[%s] expected %v got %v", name, in, out)
		}
	}
	cc, _ := http.NewClient(context.Background(),
		http.WithEndpoint(addr),
		http.WithMiddleware(mmd.Client()),
	)

	cli := &echoClient{client: NewEchoServiceHTTPClient(cc)}

	ctx := context.Background()
	ctx = metadata.NewClientContext(ctx, md)
	if out, err = cli.Echo(ctx, in); err != nil {
		t.Fatal(err)
	}
	check("echo", &SimpleMessage{Id: "test_id"}, out)

	if out, err = cli.EchoBody(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	check("echoBody", in, out)

	if out, err = cli.EchoDelete(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	check("echoDelete", in, out)

	var (
		din = &DynamicMessageUpdate{Body: &DynamicMessage{
			ValueField: &_struct.Value{Kind: &_struct.Value_StringValue{StringValue: "test"}},
		}}
		dout *DynamicMessageUpdate
	)
	if dout, err = cli.EchoResponseBody(context.Background(), din); err != nil {
		t.Fatal(err)
	}
	if din.Body.ValueField.GetStringValue() != dout.Body.ValueField.GetStringValue() {
		t.Fatalf("EchoResponseBody expected %s got %s", din, dout)
	}
	if dout, err = cli.EchoPatch(context.Background(), din); err != nil {
		t.Fatal(err)
	}
	if dout.Body == nil {
		panic("dout.body is nil")
	}
	if din.Body.ValueField.GetStringValue() != dout.Body.ValueField.GetStringValue() {
		t.Fatalf("EchoPatch expected %s got %s", din, dout)
	}
	fmt.Println("echo test success!")
}

func TestEchoGRPCServer(t *testing.T) {
	echo := &echoService{}
	ctx := context.Background()
	srv := grpc.NewServer(
		grpc.Address(":2233"),
		grpc.Middleware(mmd.Server()),
	)
	RegisterEchoServiceServer(srv, echo)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testEchoGRPCClient(t, fmt.Sprintf("127.0.0.1:2233"))
	srv.Stop(ctx)
}

func testEchoGRPCClient(t *testing.T, addr string) {
	cc, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(addr),
		grpc.WithMiddleware(mmd.Client()),
	)
	if err != nil {
		t.Fatal(err)
	}
	var (
		in  = &SimpleMessage{Id: "test_id", Num: 100}
		out = &SimpleMessage{}
	)
	client := NewEchoServiceClient(cc)
	ctx := context.Background()
	ctx = metadata.NewClientContext(ctx, md)
	if out, err = client.Echo(ctx, in); err != nil {
		t.Fatal(err)
	}
	if in.Id != out.Id || in.Num != out.Num {
		t.Errorf("expected %v got %v", in, out)
	}
}
