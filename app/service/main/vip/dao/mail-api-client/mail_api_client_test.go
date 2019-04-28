package client

import (
	"context"
	"testing"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

var cl *Client

func TestNewClient(t *testing.T) {
	var client = bm.NewClient(&bm.ClientConfig{
		App: &bm.App{
			Key:    "53e2fa226f5ad348",
			Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
		},
		Dial:      xtime.Duration(time.Second),
		Timeout:   xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second),
		Breaker: &breaker.Config{
			Window:  10 * xtime.Duration(time.Second),
			Sleep:   50 * xtime.Duration(time.Millisecond),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100,
		},
	},
	)
	convey.Convey("NewEleClient", t, func() {
		cl = NewClient(client)
		convey.So(cl, convey.ShouldNotBeNil)
	})
	convey.Convey("Get", t, func() {
		err := cl.Get(context.TODO(), "http://api.bilibili.co/x/internal/vip/user/info", nil, nil)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("Post", t, func() {
		err := cl.Post(context.TODO(), "http://api.bilibili.co/x/internal/vip/order/create", nil, nil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestClientNewClient(t *testing.T) {
	convey.Convey("NewClient", t, func(convCtx convey.C) {
		var (
			client = &bm.Client{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := NewClient(client)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestClientNewRequest(t *testing.T) {
	convey.Convey("NewRequest", t, func(convCtx convey.C) {
		var (
			method = ""
			uri    = ""
			params = interface{}(0)
			client = &bm.Client{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := NewClient(client)
			req, err := p1.NewRequest(method, uri, params)
			convCtx.Convey("Then err should be nil.req should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(req, convey.ShouldNotBeNil)
			})
		})
	})
}
