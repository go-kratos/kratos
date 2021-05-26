package testproto

import (
	context "context"
	"fmt"
	"net"
	http "net/http"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	tr "github.com/go-kratos/kratos/v2/transport/http"
	_struct "github.com/golang/protobuf/ptypes/struct"
)

type echoService struct {
}

func (s *echoService) Echo(ctx context.Context, m *SimpleMessage) (*SimpleMessage, error) {
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

func TestEchoService(t *testing.T) {
	s := &echoService{}
	h := NewEchoServiceHandler(s)
	srv := &http.Server{Addr: ":0", Handler: h}
	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	addr := lis.Addr().(*net.TCPAddr)
	time.AfterFunc(time.Second, func() {
		defer srv.Shutdown(context.Background())
		testEchoClient(t, fmt.Sprintf("127.0.0.1:%d", addr.Port))
	})
	if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func testEchoClient(t *testing.T, addr string) {
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
	cc, _ := tr.NewClient(context.Background(), tr.WithEndpoint(addr))

	cli := &echoClient{client: NewEchoServiceHTTPClient(cc)}

	if out, err = cli.Echo(context.Background(), in); err != nil {
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
}
