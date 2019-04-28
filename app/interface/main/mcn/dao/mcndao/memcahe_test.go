package mcndao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaomcnSignCacheKey(t *testing.T) {
	convey.Convey("mcnSignCacheKey", t, func(ctx convey.C) {
		var (
			mcnmid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mcnSignCacheKey(mcnmid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRawMcnSign(t *testing.T) {
	convey.Convey("RawMcnSign", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcnmid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			up, err := d.RawMcnSign(c, mcnmid)
			ctx.Convey("Then err should be nil.up should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(up, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoAsyncDelCacheMcnSign(t *testing.T) {
	convey.Convey("AsyncDelCacheMcnSign", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AsyncDelCacheMcnSign(id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaomcnDataCacheKey(t *testing.T) {
	convey.Convey("mcnDataCacheKey", t, func(ctx convey.C) {
		var (
			signID       = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mcnDataCacheKey(signID, generateDate)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoRawMcnDataSummary(t *testing.T) {
	convey.Convey("RawMcnDataSummary", t, func(ctx convey.C) {
		var (
			c            = context.Background()
			id           = int64(0)
			generateDate = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMcnDataSummary(c, id, generateDate)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}
