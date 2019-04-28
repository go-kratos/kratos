package data

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDatareverseString(t *testing.T) {
	convey.Convey("reverseString", t, func(ctx convey.C) {
		var (
			s = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := reverseString(s)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatafansRowKey(t *testing.T) {
	convey.Convey("fansRowKey", t, func(ctx convey.C) {
		var (
			id = int64(0)
			ty = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := fansRowKey(id, ty)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataarcQueryKey(t *testing.T) {
	convey.Convey("arcQueryKey", t, func(ctx convey.C) {
		var (
			id = int64(0)
			dt = ""
			cp = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := arcQueryKey(id, dt, cp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataUpFansAnalysis(t *testing.T) {
	convey.Convey("UpFansAnalysis", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ty  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpFansAnalysis(c, mid, ty)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataUpArcQuery(t *testing.T) {
	convey.Convey("UpArcQuery", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = ""
			cp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpArcQuery(c, mid, dt, cp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}
