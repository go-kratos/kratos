package assist

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistassistKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("assistKey", t, func(ctx convey.C) {
		p1 := assistKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistSetCacheAss(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(0)
		assistMids = []int64{}
	)
	convey.Convey("SetCacheAss", t, func(ctx convey.C) {
		err := d.SetCacheAss(c, mid, assistMids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistGetCacheAss(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("GetCacheAss", t, func(ctx convey.C) {
		assistMids, err := d.GetCacheAss(c, mid)
		ctx.Convey("Then err should be nil.assistMids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assistMids, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistDelCacheAss(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelCacheAss", t, func(ctx convey.C) {
		err := d.DelCacheAss(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistpingMemcache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingMemcache", t, func(ctx convey.C) {
		err := d.pingMemcache(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
