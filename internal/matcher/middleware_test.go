package matcher

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/middleware"
)

func logging(module string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			return module, nil
		}
	}
}

func equal(ms []middleware.Middleware, modules ...string) bool {
	if len(ms) == 0 {
		return false
	}
	for i, m := range ms {
		x, _ := m(nil)(nil, nil)
		if x != modules[i] {
			return false
		}
	}
	return true
}

func TestMatcher(t *testing.T) {
	m := New()
	m.Use(logging("logging"))
	m.Add("/*", logging("*"))
	m.Add("/foo/bar/*", logging("foo/bar/*"))
	m.Add("/foo/bar", logging("foo/bar"))
	m.Add("/foo/*", logging("foo/*"))

	if ms := m.Match("/"); len(ms) != 2 {
		t.Fatal("not equal")
	} else if !equal(ms, "logging", "*") {
		t.Fatal("not equal")
	}

	if ms := m.Match("/foo/xxx"); len(ms) != 3 {
		t.Fatal("not equal")
	} else if !equal(ms, "logging", "*", "foo/*") {
		t.Fatal("not equal")
	}

	if ms := m.Match("/foo/bar"); len(ms) != 4 {
		t.Fatal("not equal")
	} else if !equal(ms, "logging", "*", "foo/*", "foo/bar") {
		t.Fatal("not equal")
	}

	if ms := m.Match("/foo/bar/x"); len(ms) != 4 {
		t.Fatal("not equal")
	} else if !equal(ms, "logging", "*", "foo/*", "foo/bar/*") {
		t.Fatal("not equal")
	}
}
