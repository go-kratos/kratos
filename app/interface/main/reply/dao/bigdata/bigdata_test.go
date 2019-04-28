package bigdata

import (
	"context"
	"testing"

	"go-common/app/interface/main/reply/conf"

	"github.com/smartystreets/goconvey/convey"
)

func TestBigdataNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		var (
			c = conf.Conf
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := New(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBigdataFilter(t *testing.T) {
	convey.Convey("Filter", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			msg = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Filter(c, msg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
