package sentry

import (
	"context"
	"testing"
	"time"
)

func TestWithTags(t *testing.T) {
	opts := new(options)
	strval := "bar"
	kvs := map[string]string{
		"foo": strval,
	}
	funcTags := WithTags(kvs)
	funcTags(opts)
	if opts.tags["foo"] != strval {
		t.Errorf("TestWithTags() = %v, want %v", opts.tags["foo"], strval)
	}
}

func TestWithContextTags(t *testing.T) {
	opts := new(options)
	fn := func(context.Context) map[string]string {
		return map[string]string{"foo": "bar"}
	}
	funcTags := WithContextTags(fn)
	funcTags(opts)
	if opts.contextTags == nil {
		t.Fatal("contextTags is nil")
	}
	if got := opts.contextTags(context.Background())["foo"]; got != "bar" {
		t.Errorf("TestWithContextTags() = %v, want %v", got, "bar")
	}
}

func TestWithRepanic(t *testing.T) {
	opts := new(options)
	val := true
	f := WithRepanic(val)
	f(opts)
	if opts.repanic != val {
		t.Errorf("TestWithRepanic() = %v, want %v", opts.repanic, val)
	}
}

func TestWithWaitForDelivery(t *testing.T) {
	opts := new(options)
	val := true
	f := WithWaitForDelivery(val)
	f(opts)
	if opts.waitForDelivery != val {
		t.Errorf("TestWithWaitForDelivery() = %v, want %v", opts.waitForDelivery, val)
	}
}

func TestWithTimeout(t *testing.T) {
	opts := new(options)
	val := time.Second * 10
	f := WithTimeout(val)
	f(opts)
	if opts.timeout != val {
		t.Errorf("TestWithTimeout() = %v, want %v", opts.timeout, val)
	}
}
