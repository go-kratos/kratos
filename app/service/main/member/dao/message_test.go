package dao

import (
	"context"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"gopkg.in/h2non/gock.v1"
)

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestDaoSendMessage(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(123)
		title = "test"
		msg   = "test"
		mc    = "test"
	)
	convey.Convey("SendMessage", t, func(ctx convey.C) {
		d.client.SetTransport(gock.DefaultTransport)
		defer gock.OffAll()
		httpMock("POST", notifyURL).Reply(200).JSON(`{"code":0,"message":"0"}`)
		err := d.SendMessage(c, mid, title, msg, mc)
		ctx.So(err, convey.ShouldBeNil)
	})
}
