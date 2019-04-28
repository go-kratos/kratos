package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveUppersCount(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{1, 2}
	)
	convey.Convey("UppersCount", t, func(ctx convey.C) {
		uc, err := d.UppersCount(c, mids)
		ctx.Convey("Then err should be nil.uc should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(uc, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpperCount(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("UpperCount", t, func(ctx convey.C) {
		count, err := d.UpperCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpperPassed(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("UpperPassed", t, func(ctx convey.C) {
		_, _, _, err := d.UpperPassed(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveUppersPassed(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{1}
	)
	convey.Convey("UppersPassed", t, func(ctx convey.C) {
		aidm, ptimes, copyrights, err := d.UppersPassed(c, mids)
		ctx.Convey("Then err should be nil.aidm,ptimes,copyrights should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(copyrights, convey.ShouldNotBeNil)
			ctx.So(ptimes, convey.ShouldNotBeNil)
			ctx.So(aidm, convey.ShouldNotBeNil)
		})
	})
}
