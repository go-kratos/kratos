package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccTags(t *testing.T) {
	convey.Convey("AccTags", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.AccTags(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsAnswered(t *testing.T) {
	convey.Convey("IsAnswered", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2222)
			start = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			status, err := d.IsAnswered(c, mid, start)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}
