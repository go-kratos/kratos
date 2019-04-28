package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_Bnj2019Conf(t *testing.T) {
	convey.Convey("ElecShow", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rs, err := d.Bnj2019Conf(c)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", rs)
			})
		})
	})
}
