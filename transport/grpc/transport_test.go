package grpc

import (
	"reflect"
	"sort"
	"testing"

	"github.com/go-kratos/kratos/v2/transport"
)

func TestTransport_Kind(t *testing.T) {
	o := &Transport{}
	if !reflect.DeepEqual(transport.KindGRPC, o.Kind()) {
		t.Errorf("expect %v, got %v", transport.KindGRPC, o.Kind())
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

func TestTransport_RequestHeader(t *testing.T) {
	v := headerCarrier{}
	v.Set("a", "1")
	o := &Transport{reqHeader: v}
	if !reflect.DeepEqual("1", o.RequestHeader().Get("a")) {
		t.Errorf("expect %v, got %v", "1", o.RequestHeader().Get("a"))
	}
	if !reflect.DeepEqual("", o.RequestHeader().Get("notfound")) {
		t.Errorf("expect %v, got %v", "", o.RequestHeader().Get("notfound"))
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

func TestHeaderCarrier_Keys(t *testing.T) {
	v := headerCarrier{}
	v.Set("abb", "1")
	v.Set("bcc", "2")
	want := []string{"abb", "bcc"}
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
