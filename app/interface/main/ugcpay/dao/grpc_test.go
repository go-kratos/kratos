package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTradeCreate(t *testing.T) {
	convey.Convey("TradeCreate", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			platform = "ios"
			mid      = int64(46333)
			oid      = int64(10110745)
			otype    = "archive"
			currency = "bp"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			orderID, payData, err := d.TradeCreate(c, platform, mid, oid, otype, currency)
			ctx.Convey("Then err should be nil.orderID,payData should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(payData, convey.ShouldNotBeNil)
				ctx.So(orderID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTradeQuery(t *testing.T) {
	convey.Convey("TradeQuery", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "32172647181109193324"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			order, err := d.TradeQuery(c, orderID)
			ctx.Convey("Then err should be nil.order should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(order, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTradeConfirm(t *testing.T) {
	convey.Convey("TradeConfirm", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "32172647181109193324"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			order, err := d.TradeConfirm(c, orderID)
			ctx.Convey("Then err should be nil.order should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(order, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTradeCancel(t *testing.T) {
	convey.Convey("TradeCancel", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = "32172647181109193324"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TradeCancel(c, orderID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIncomeAssetOverview(t *testing.T) {
	convey.Convey("IncomeAssetOverview", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			inc, err := d.IncomeAssetOverview(c, mid)
			ctx.Convey("Then err should be nil.inc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(inc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIncomeUserAssetList(t *testing.T) {
	convey.Convey("IncomeUserAssetList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(46333)
			ver = int64(201811)
			ps  = int64(10)
			pn  = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			inc, err := d.IncomeUserAssetList(c, mid, ver, ps, pn)
			ctx.Convey("Then err should be nil.inc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(inc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveTitles(t *testing.T) {
	convey.Convey("ArchiveTitles", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{10110745}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcTitles, err := d.ArchiveTitles(c, aids)
			ctx.Convey("Then err should be nil.arcTitles should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(arcTitles, convey.ShouldNotBeNil)
			})
		})
	})
}
