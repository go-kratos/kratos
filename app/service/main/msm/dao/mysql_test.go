package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCodes(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Codes", t, func(ctx convey.C) {
		codes, lcode, err := d.Codes(c)
		ctx.Convey("Then err should be nil.codes,lcode should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lcode, convey.ShouldNotBeNil)
			ctx.So(codes, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDiff(t *testing.T) {
	var (
		c   = context.Background()
		ver = int64(0)
	)
	convey.Convey("Diff", t, func(ctx convey.C) {
		vers, err := d.Diff(c, ver)
		ctx.Convey("Then err should be nil.vers should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vers, convey.ShouldNotBeNil)
		})
	})
}
