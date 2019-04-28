package search

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchSetSearchAuditStat(t *testing.T) {
	convey.Convey("SetSearchAuditStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			key   = "test"
			state bool
		)
		state = true
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetSearchAuditStat(c, key, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestSearchGetSearchAuditStat(t *testing.T) {
	convey.Convey("GetSearchAuditStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			f, date, err := d.GetSearchAuditStat(c, key)
			ctx.Convey("Then err should be nil.f,date should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(date, convey.ShouldNotBeNil)
				ctx.So(f, convey.ShouldNotBeNil)
			})
		})
	})
}
