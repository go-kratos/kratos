package sentry

import (
	"testing"
	"time"
)

func TestWithTags(t *testing.T) {
	opts := new(options)
	strval := "bar"
	kvs := map[string]interface{}{
		"foo": strval,
	}
	funcTags := WithTags(kvs)
	funcTags(opts)
	if opts.tags["foo"].(string) != strval {
		t.Errorf("TestWithTags() = %v, want %v", opts.tags["foo"].(string), strval)
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
