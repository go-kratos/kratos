package pgc

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPgcCardsInfoReply(t *testing.T) {
	convey.Convey("CardsInfoReply", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			seasonIds = []int32{33730}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CardsInfoReply(c, seasonIds)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPgcCardsEpInfoReply(t *testing.T) {
	convey.Convey("CardsEpInfoReply", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			epIds = []int32{117117}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CardsEpInfoReply(c, epIds)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
