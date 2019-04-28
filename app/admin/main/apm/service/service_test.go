package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService(t *testing.T) {
	Convey("service", t, func() {
		t.Log("service test")
	})
}

func TestPing(t *testing.T) {
	Convey("ping", t, func() {
		var s = &Service{}
		s.Ping(context.Background())
	})
}
