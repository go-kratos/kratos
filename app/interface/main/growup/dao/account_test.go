package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccountInfos(t *testing.T) {
	convey.Convey("AccountInfos", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1001}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.AccountInfos(c, mids)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpBusinessInfos(t *testing.T) {
	convey.Convey("UpBusinessInfos", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpBusinessInfos(c, mid)
			ctx.Convey("Then err should be nil.identify should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCard(t *testing.T) {
	convey.Convey("Card", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Card(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoProfileWithStat(t *testing.T) {
	convey.Convey("ProfileWithStat", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ProfileWithStat(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
