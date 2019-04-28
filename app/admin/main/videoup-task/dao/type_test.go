package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTypeMapping(t *testing.T) {
	convey.Convey("TypeMapping", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tmap, err := d.TypeMapping(c)
			ctx.Convey("Then err should be nil.tmap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tmap, convey.ShouldNotBeNil)
			})
		})
	})
}
