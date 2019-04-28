package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaotopDyHit(t *testing.T) {
	convey.Convey("topDyHit", t, func(ctx convey.C) {
		var (
			mid = int64(666)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := topDyHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaotopDyKey(t *testing.T) {
	convey.Convey("topDyKey", t, func(ctx convey.C) {
		var (
			mid = int64(666)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := topDyKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawTopDynamic(t *testing.T) {
	convey.Convey("RawTopDynamic", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(666)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			dyID, err := d.RawTopDynamic(c, mid)
			ctx.Convey("Then err should be nil.dyID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(dyID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddTopDynamic(t *testing.T) {
	convey.Convey("AddTopDynamic", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(666)
			dyID = int64(1111111111)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddTopDynamic(c, mid, dyID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelTopDynamic(t *testing.T) {
	convey.Convey("DelTopDynamic", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(666)
			now = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelTopDynamic(c, mid, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
