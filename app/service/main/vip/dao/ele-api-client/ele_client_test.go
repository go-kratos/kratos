package client

import (
	"bytes"
	"context"
	"testing"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

var cl *EleClient

func TestNewEleClient(t *testing.T) {
	var c = &Config{
		App: &App{
			Key:    "sdfsdf",
			Secret: "sdfsdf",
		},
	}
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
		cl = NewEleClient(c, client)
		convey.So(cl, convey.ShouldNotBeNil)
	})
	convey.Convey("Get", t, func() {
		err := cl.Get(context.TODO(), "http://api.bilibili.co", "/x/internal/vip/user/info", nil, nil)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("Post", t, func() {
		err := cl.Post(context.TODO(), "http://api.bilibili.co", "/x/internal/vip/order/create", nil, nil)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("newRequest", t, func() {
		req, err := cl.newRequest("POST", "http://api.bilibili.co", "/x/internal/vip/user/info", nil)
		convey.So(err, convey.ShouldBeNil)
		convey.So(req, convey.ShouldNotBeNil)
	})
}

func TestIsSuccess(t *testing.T) {
	convey.Convey("IsSuccess", t, func() {
		p1 := IsSuccess("ok")
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestEleSign(t *testing.T) {
	convey.Convey("eleSign", t, func() {
		p1 := eleSign("", "", "", "", "")
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestComputeHmac256(t *testing.T) {
	convey.Convey("computeHmac256", t, func() {
		var b bytes.Buffer
		b.WriteString("http://bilibili.com/x/vip")
		b.WriteString("&")
		b.WriteString("consumer_key=")

		p1 := computeHmac256(b, "xxx")
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestUUID4(t *testing.T) {
	convey.Convey("UUID4", t, func() {
		p1 := UUID4()
		convey.So(p1, convey.ShouldNotBeNil)
	})
}
