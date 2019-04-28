package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopingMC(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoSetRegionCache(t *testing.T) {
	var (
		c          = context.Background()
		regionArcs map[int32][]int64
	)
	convey.Convey("SetRegionCache", t, func(ctx convey.C) {
		err := d.SetRegionCache(c, regionArcs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoRegionCache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("RegionCache", t, func(ctx convey.C) {
		rs := d.RegionCache(c)
		ctx.Convey("Then rs should not be nil.", func(ctx convey.C) {
			ctx.So(rs, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetTagCache(t *testing.T) {
	var (
		c             = context.Background()
		regionTagArcs map[string][]int64
	)
	convey.Convey("SetTagCache", t, func(ctx convey.C) {
		err := d.SetTagCache(c, regionTagArcs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTagCache(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("TagCache", t, func(ctx convey.C) {
		rs := d.TagCache(c)
		ctx.Convey("Then rs should not be nil.", func(ctx convey.C) {
			ctx.So(rs, convey.ShouldNotBeNil)
		})
	})
}
