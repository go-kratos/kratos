package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocouponBalancesKey(t *testing.T) {
	convey.Convey("couponBalancesKey", t, func(convCtx convey.C) {
		var (
			mid = int64(1)
			ct  = int8(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := couponBalancesKey(mid, ct)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldEqual, "cbl:2:1")
			})
		})
	})
}

func TestDaoDelCouponBalancesCache(t *testing.T) {
	convey.Convey("DelCouponBalancesCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ct  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCouponBalancesCache(c, mid, ct)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaocouponsKey(t *testing.T) {
	convey.Convey("couponsKey", t, func(convCtx convey.C) {
		var (
			mid = int64(22)
			ct  = int8(33)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := couponsKey(mid, ct)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldEqual, "cs:33:22")
			})
		})
	})
}

func TestDaoprizeCardKey(t *testing.T) {
	convey.Convey("prizeCardKey", t, func(convCtx convey.C) {
		var (
			mid      = int64(22)
			actID    = int64(1)
			cardType = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := prizeCardKey(mid, actID, cardType)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldEqual, "nypc:22:1:0")
			})
		})
	})
}

func TestDaoprizeCardsKey(t *testing.T) {
	convey.Convey("prizeCardsKey", t, func(convCtx convey.C) {
		var (
			mid   = int64(33)
			actID = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := prizeCardsKey(mid, actID)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldEqual, "nypcs:33:1")
			})
		})
	})
}

func TestDaoDelCouponsCache(t *testing.T) {
	convey.Convey("DelCouponsCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
			ct  = int8(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCouponsCache(c, mid, ct)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaodelCache(t *testing.T) {
	convey.Convey("delCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			key = "1"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.delCache(c, key)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
			err = d.delCache(c, "")
			convCtx.Convey("Then err should be not nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPrizeCardKey(t *testing.T) {
	convey.Convey("DelPrizeCardKey", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			actID    = int64(0)
			cardType = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelPrizeCardKey(c, mid, actID, cardType)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelPrizeCardsKey(t *testing.T) {
	convey.Convey("DelPrizeCardsKey", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			actID = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelPrizeCardsKey(c, mid, actID)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
