package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBlacklist(t *testing.T) {
	convey.Convey("Blacklist", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			blacklist, err := d.Blacklist(c)
			ctx.Convey("Then err should be nil.blacklist should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(blacklist, convey.ShouldNotBeNil)
			})
		})
	})
}
