package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMidsByName(t *testing.T) {
	var (
		c     = context.TODO()
		names = []string{"", "一"}
	)
	convey.Convey("Search mids by name", t, func(ctx convey.C) {
		p1, err := d.MidsByName(c, names)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
