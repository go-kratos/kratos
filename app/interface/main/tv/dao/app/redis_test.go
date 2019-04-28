package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyZone(t *testing.T) {
	var category = int(0)
	convey.Convey("keyZone", t, func(c convey.C) {
		p1 := keyZone(category)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoZrevrangeList(t *testing.T) {
	var (
		c        = context.Background()
		category = int(1)
		start    = int(0)
		end      = int(10)
	)
	convey.Convey("ZrevrangeList", t, func(ctx convey.C) {
		sids, count, err := d.ZrevrangeList(c, category, start, end)
		ctx.Convey("Then err should be nil.sids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(sids, convey.ShouldNotBeNil)
		})
	})
}
