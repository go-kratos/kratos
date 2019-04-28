package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocouponAllowancesKey(t *testing.T) {
	convey.Convey("couponAllowancesKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := couponAllowancesKey(mid, 0)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocouponsKey(t *testing.T) {
	convey.Convey("couponsKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
			ct  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := couponsKey(mid, ct)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaograntUnique(t *testing.T) {
	convey.Convey("grantUnique", t, func(convCtx convey.C) {
		var (
			token = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := grantUnique(token)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelCouponAllowancesKey(t *testing.T) {
	convey.Convey("DelCouponAllowancesKey", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCouponAllowancesKey(c, mid, 0)
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
			key = "123"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.delCache(c, key)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCouponTypeCache(t *testing.T) {
	convey.Convey("DelCouponTypeCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ct  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelCouponTypeCache(c, mid, ct)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelGrantUniqueLock(t *testing.T) {
	convey.Convey("DelGrantUniqueLock", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			token = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelGrantUniqueLock(c, token)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddGrantUniqueLock(t *testing.T) {
	convey.Convey("AddGrantUniqueLock", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			token   = ""
			seconds = int32(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			succeed := d.AddGrantUniqueLock(c, token, seconds)
			convCtx.Convey("Then succeed should not be nil.", func(convCtx convey.C) {
				convCtx.So(succeed, convey.ShouldNotBeNil)
			})
		})
	})
}
