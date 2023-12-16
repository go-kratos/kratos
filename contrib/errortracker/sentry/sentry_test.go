package sentry

import "testing"

func TestWithTags(t *testing.T) {
	opts := new(options)
	strval := "bar"
	kvs := map[string]interface{}{
		"foo": strval,
	}
	funcTags := WithTags(kvs)
	funcTags(opts)
	if opts.tags["foo"].(string) != strval {
		t.Errorf("TestWithTags() = %s, want %s", opts.tags["foo"].(string), strval)
	}
}
