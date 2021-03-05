package testproto

import (
	"bytes"
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
	baseURL string
	client  *http.Client
}

// post: /v1/example/echo/{id}
func (c *echoClient) Echo(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(in)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/example/echo/%s", c.baseURL, in.Id), bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	out = new(SimpleMessage)
	if err = tr.Do(c.client, req, out); err != nil {
		return
	}
	return
}

// post: /v1/example/echo_body
func (c *echoClient) EchoBody(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(in)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/example/echo_body", c.baseURL), bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	out = new(SimpleMessage)
	if err = tr.Do(c.client, req, out); err != nil {
		return
	}
	return
}

// delete: /v1/example/echo_delete/{id}/{num}
func (c *echoClient) EchoDelete(ctx context.Context, in *SimpleMessage) (out *SimpleMessage, err error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/example/echo_delete/%s/%d", c.baseURL, in.Id, in.Num), nil)
	if err != nil {
		return
	}
	out = new(SimpleMessage)
	if err = tr.Do(c.client, req, out); err != nil {
		return
	}
	return
}

// patch: /v1/example/echo_patch
func (c *echoClient) EchoPatch(ctx context.Context, in *DynamicMessageUpdate) (out *DynamicMessageUpdate, err error) {
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(in.Body)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/v1/example/echo_patch", c.baseURL), bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	out = new(DynamicMessageUpdate)
	if err = tr.Do(c.client, req, out); err != nil {
		return
	}
	return
}

// post: /v1/example/echo_response_body
func (c *echoClient) EchoResponseBody(ctx context.Context, in *DynamicMessageUpdate) (out *DynamicMessageUpdate, err error) {
	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(in)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/example/echo_response_body", c.baseURL), bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	out = new(DynamicMessageUpdate)
	if err = tr.Do(c.client, req, out.Body); err != nil {
		return
	}
	return
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
		testEchoClient(t, fmt.Sprintf("http://127.0.0.1:%d", addr.Port))
	})
	if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func testEchoClient(t *testing.T, baseURL string) {
	var (
		err error
		in  = &SimpleMessage{Id: "test_id", Num: 100}
		out = &SimpleMessage{}
	)
	check := func(in, out *SimpleMessage) {
		if in.Id != out.Id || in.Num != out.Num {
			t.Errorf("expected %s got %s", in, out)
		}
	}

	cli := &echoClient{baseURL: baseURL, client: http.DefaultClient}

	if out, err = cli.Echo(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	check(in, out)

	if out, err = cli.EchoBody(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	check(in, out)

	if out, err = cli.EchoDelete(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	check(in, out)

	var (
		din = &DynamicMessageUpdate{Body: &DynamicMessage{
			ValueField: &_struct.Value{Kind: &_struct.Value_StringValue{StringValue: "test"}},
		}}
		dout *DynamicMessageUpdate
	)
	if dout, err = cli.EchoPatch(context.Background(), din); err != nil {
		t.Fatal(err)
	}
	if din.Body.ValueField.GetStringValue() != dout.Body.ValueField.GetStringValue() {
		t.Errorf("EchoPatch expected %s got %s", din, dout)
	}
	if dout, err = cli.EchoResponseBody(context.Background(), din); err != nil {
		t.Fatal(err)
	}
	if din.Body.ValueField.GetStringValue() != dout.Body.ValueField.GetStringValue() {
		t.Errorf("EchoResponseBody expected %s got %s", din, dout)
	}
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
