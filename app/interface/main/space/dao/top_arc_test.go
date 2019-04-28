package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotopArcHit(t *testing.T) {
	convey.Convey("topArcHit", t, func(ctx convey.C) {
		var (
			mid = int64(5551)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := topArcHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotopArcKey(t *testing.T) {
	convey.Convey("topArcKey", t, func(ctx convey.C) {
		var (
			mid = int64(5551)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := topArcKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawTopArc(t *testing.T) {
	convey.Convey("RawTopArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.RawTopArc(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddTopArc(t *testing.T) {
	convey.Convey("AddTopArc", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(5551)
			aid    = int64(170001)
			reason = "aaaa"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddTopArc(c, mid, aid, reason)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelTopArc(t *testing.T) {
	convey.Convey("DelTopArc", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(5551)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelTopArc(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
