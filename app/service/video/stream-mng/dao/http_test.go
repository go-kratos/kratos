package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNewRequst(t *testing.T) {
	convey.Convey("NewRequst", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			method = ""
			url    = "http://live-stream.bilibili.co/"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.NewRequst(c, method, url, nil, nil, nil, nil)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
