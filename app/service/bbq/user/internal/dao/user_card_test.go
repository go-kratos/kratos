package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawUserCard(t *testing.T) {
	convey.Convey("RawUserCard", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			userCard, err := d.RawUserCard(c, mid)
			ctx.Convey("Then err should be nil.userCard should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(userCard, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawUserCards(t *testing.T) {
	convey.Convey("RawUserCards", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			userCards, err := d.RawUserCards(c, mids)
			ctx.Convey("Then err should be nil.userCards should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(userCards, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawUserAccCards(t *testing.T) {
	convey.Convey("RawUserAccCards", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawUserAccCards(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
