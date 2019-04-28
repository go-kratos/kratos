package creative

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBgmData(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("BgmData", t, func(ctx convey.C) {
		_, err := d.RawBgmData(c, 1, 1, 1)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCacheBgmData1(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("BgmData", t, func(ctx convey.C) {
		_, err := d.BgmData(c, 1, 1, 1, false)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAddCacheBgmData(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("AddCacheBgmData", t, func(ctx convey.C) {
		err := d.AddCacheBgmData(c, 1, 1, 1, nil)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDelCacheBgmData(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("DelCacheBgmData", t, func(ctx convey.C) {
		err := d.DelCacheBgmData(c, 1, 1, 1)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCacheBgmData(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("CacheBgmData", t, func(ctx convey.C) {
		_, err := d.CacheBgmData(c, 1, 1, 1)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
