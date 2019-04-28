package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccountCards(t *testing.T) {
	convey.Convey("AccountCards", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{46333}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			cards, err := d.AccountCards(c, mids)
			ctx.Convey("Then err should be nil.cards should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cards, convey.ShouldHaveLength, 1)
			})
		})
	})
}
