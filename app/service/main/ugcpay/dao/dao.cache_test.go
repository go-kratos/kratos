package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOrderUser(t *testing.T) {
	convey.Convey("OrderUser", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = "56272730181109202542"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.OrderUser(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAsset(t *testing.T) {
	convey.Convey("Asset", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(10110745)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelCacheAsset(c, id, otype, currency)
			convey.So(err, convey.ShouldBeNil)
			res, err := d.Asset(c, id, otype, currency)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

// func TestDaoAggrIncomeUser(t *testing.T) {
// 	convey.Convey("AggrIncomeUser", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			mid = int64(46333)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			res, err := d.AggrIncomeUser(c, mid, cur)
// 			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 				ctx.So(res, convey.ShouldNotBeNil)
// 			})
// 		})
// 	})
// }

// func TestDaoAggrIncomeUserMonthly(t *testing.T) {
// 	convey.Convey("AggrIncomeUserMonthly", t, func(ctx convey.C) {
// 		var (
// 			c   = context.Background()
// 			id  = int64(46333)
// 			ver = int64(201810)
// 			cur = "bp"
// 		)
// 		ctx.Convey("When everything goes positive", func(ctx convey.C) {
// 			res, err := d.AggrIncomeUseAssetList(c, id, cur, ver)
// 			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
// 				ctx.So(err, convey.ShouldBeNil)
// 				ctx.So(res, convey.ShouldNotBeNil)
// 			})
// 		})
// 	})
// }
