package client

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func startRPCServer() {
	s = New(nil)
}

func TestClient_DelTokenCache(t *testing.T) {
	startRPCServer()
	time.Sleep(3 * time.Second)
	Convey("Test RPC client del token by token", t, func() {
		var (
			c  = context.TODO()
			ak = "64294c76972aee8cf4af51566c76ed0d"
		)
		err := s.DelTokenCache(c, ak)
		So(err, ShouldBeNil)
	})
}

func TestClient_DelCookieCookie(t *testing.T) {
	startRPCServer()
	time.Sleep(3 * time.Second)
	Convey("Test RPC Client del cookie by cookie", t, func() {
		sd := "1d0fb9cf,1,7f9745b6"
		err := s.DelCookieCookie(context.TODO(), sd)
		So(err, ShouldBeNil)
	})
}
