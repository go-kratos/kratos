package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocacheSFstreamFullInfo(t *testing.T) {
	convey.Convey("cacheSFstreamFullInfo", t, func(ctx convey.C) {
		var (
			id    = int64(11891462)
			sname = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.cacheSFstreamFullInfo(id, sname)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocacheSFstreamRIDByName(t *testing.T) {
	convey.Convey("cacheSFstreamRIDByName", t, func(ctx convey.C) {
		var (
			sname = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.cacheSFstreamRIDByName(sname)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStreamFullInfo(t *testing.T) {
	convey.Convey("StreamFullInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int64(11891462)
			sname = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.StreamFullInfo(c, rid, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOriginUpStreamInfo(t *testing.T) {
	convey.Convey("OriginUpStreamInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sname, origin, err := d.OriginUpStreamInfo(c, rid)
			ctx.Convey("Then err should be nil.sname,origin should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(origin, convey.ShouldNotBeNil)
				ctx.So(sname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOriginUpStreamInfoBySName(t *testing.T) {
	convey.Convey("OriginUpStreamInfoBySName", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sname = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rid, origin, err := d.OriginUpStreamInfoBySName(c, sname)
			ctx.Convey("Then err should be nil.rid,origin should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(origin, convey.ShouldNotBeNil)
				ctx.So(rid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStreamRIDByName(t *testing.T) {
	convey.Convey("StreamRIDByName", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sname = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.StreamRIDByName(c, sname)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMultiStreamInfo(t *testing.T) {
	convey.Convey("MultiStreamInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{11891462}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MultiStreamInfo(c, rids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
