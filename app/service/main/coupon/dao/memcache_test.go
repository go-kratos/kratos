package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/coupon/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoreceiveLogKey(t *testing.T) {
	var (
		appkey  = "123"
		orderNo = "456"
		ct      = int8(0)
	)
	convey.Convey("TestDaoreceiveLogKey ", t, func(ctx convey.C) {
		p1 := receiveLogKey(appkey, orderNo, ct)
		ctx.Convey("Then p1 should equal.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "rl:1234560")
		})
	})
}

func TestDaoprizeCardKey(t *testing.T) {
	var (
		mid   int64 = 22
		actID int64 = 1
		ct          = int8(0)
	)
	convey.Convey("TestDaoprizeCardKey ", t, func(ctx convey.C) {
		p1 := prizeCardKey(mid, actID, ct)
		ctx.Convey("Then p1 should equal.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "nypc:22:1:0")
		})
	})
}

func TestDaoprizeCardsKey(t *testing.T) {
	var (
		mid   int64 = 22
		actID int64 = 1
	)
	convey.Convey("TestDaoprizeCardsKey ", t, func(ctx convey.C) {
		p1 := prizeCardsKey(mid, actID)
		ctx.Convey("Then p1 should equal.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "nypcs:22:1")
		})
	})
}

func TestDaocouponuniqueNoKey(t *testing.T) {
	var (
		uniqueno string = "uniqueno"
	)
	convey.Convey("TestDaocouponuniqueNoKey ", t, func(ctx convey.C) {
		p1 := couponuniqueNoKey(uniqueno)
		ctx.Convey("Then p1 should equal.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldEqual, "uq:uniqueno")
		})
	})
}
func TestDaocouponsKey(t *testing.T) {
	var (
		mid = int64(0)
		ct  = int8(0)
	)
	convey.Convey("couponsKey ", t, func(ctx convey.C) {
		p1 := couponsKey(mid, ct)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaouseUniqueKey(t *testing.T) {
	var (
		orderNO = "1"
		ct      = int8(0)
	)
	convey.Convey("useUniqueKey", t, func(ctx convey.C) {
		p1 := useUniqueKey(orderNO, ct)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaocouponBalancesKey(t *testing.T) {
	var (
		mid = int64(0)
		ct  = int8(0)
	)
	convey.Convey("couponBalancesKey", t, func(ctx convey.C) {
		p1 := couponBalancesKey(mid, ct)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaouserGrantKey(t *testing.T) {
	var (
		token = "1"
		mid   = int64(0)
	)
	convey.Convey("userGrantKey", t, func(ctx convey.C) {
		p1 := userGrantKey(token, mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaobranchCurrentCount(t *testing.T) {
	var (
		token = "1"
	)
	convey.Convey("branchCurrentCount", t, func(ctx convey.C) {
		p1 := branchCurrentCount(token)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaocouponAllowancesKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("couponAllowancesKey", t, func(ctx convey.C) {
		p1 := couponAllowancesKey(mid, 0)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelUniqueKey(t *testing.T) {
	var (
		c       = context.TODO()
		orderNO = "1"
		ct      = int8(0)
	)
	convey.Convey("DelUniqueKey", t, func(ctx convey.C) {
		err := d.DelUniqueKey(c, orderNO, ct)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCouponsCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ct  = int8(0)
	)
	convey.Convey("DelCouponsCache", t, func(ctx convey.C) {
		err := d.DelCouponsCache(c, mid, ct)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCouponBalancesCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ct  = int8(0)
	)
	convey.Convey("DelCouponBalancesCache", t, func(ctx convey.C) {
		err := d.DelCouponBalancesCache(c, mid, ct)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelGrantKey(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
		mid   = int64(0)
	)
	convey.Convey("DelGrantKey", t, func(ctx convey.C) {
		err := d.DelGrantKey(c, token, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelBranchCurrentCountKey(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
	)
	convey.Convey("DelBranchCurrentCountKey", t, func(ctx convey.C) {
		err := d.DelBranchCurrentCountKey(c, token)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCouponAllowancesKey(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(0)
	)
	convey.Convey("DelCouponAllowancesKey", t, func(ctx convey.C) {
		err := d.DelCouponAllowancesKey(c, mid, 0)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelPrizeCardKey(t *testing.T) {
	var (
		c           = context.Background()
		mid   int64 = 22
		actID int64 = 1
		ct          = int8(0)
	)
	convey.Convey("DelPrizeCardKey", t, func(ctx convey.C) {
		err := d.DelPrizeCardKey(c, mid, actID, ct)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelPrizeCardsKey(t *testing.T) {
	var (
		c           = context.Background()
		mid   int64 = 22
		actID int64 = 1
	)
	convey.Convey("DelCouponAllowancesKey", t, func(ctx convey.C) {
		err := d.DelPrizeCardsKey(c, mid, actID)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCouponsCache(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		ct      = int8(0)
		coupons = []*model.CouponInfo{}
		err     error
	)
	convey.Convey("CouponsCache", t, func(ctx convey.C) {
		coupons, err = d.CouponsCache(c, mid, ct)
		ctx.Convey("Then err should be nil.coupons should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(coupons, convey.ShouldBeNil)
		})
		err = d.SetCouponsCache(c, mid, ct, coupons)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		coupons, err = d.CouponsCache(c, mid, ct)
		ctx.Convey("Then err should be nil.coupons should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(coupons, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddUseUniqueLock(t *testing.T) {
	var (
		c       = context.TODO()
		orderNO = "1"
		ct      = int8(0)
	)
	convey.Convey("AddUseUniqueLock", t, func(ctx convey.C) {
		succeed := d.AddUseUniqueLock(c, orderNO, ct)
		ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
			ctx.So(succeed, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddReceiveUniqueLock(t *testing.T) {
	var (
		c       = context.TODO()
		appkey  = "1"
		orderNO = "2"
		ct      = int8(0)
	)
	convey.Convey("AddReceiveUniqueLock", t, func(ctx convey.C) {
		succeed := d.AddReceiveUniqueLock(c, appkey, orderNO, ct)
		ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
			ctx.So(succeed, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelReceiveUniqueLock(t *testing.T) {
	var (
		c       = context.TODO()
		appkey  = "1"
		orderNO = "1"
		ct      = int8(0)
	)
	convey.Convey("DelReceiveUniqueLock ", t, func(ctx convey.C) {
		err := d.DelReceiveUniqueLock(c, appkey, orderNO, ct)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaodelCache(t *testing.T) {
	var (
		c   = context.TODO()
		key = "1"
	)
	convey.Convey("delCache", t, func(ctx convey.C) {
		err := d.delCache(c, key)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCouponBlanceCache(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		ct      = int8(0)
		coupons = []*model.CouponBalanceInfo{}
	)
	convey.Convey("CouponBlanceCache", t, func(ctx convey.C) {
		err := d.SetCouponBlanceCache(c, mid, ct, coupons)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		coupons, err := d.CouponBlanceCache(c, mid, ct)
		ctx.Convey("Then err should be nil.coupons should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(coupons, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddUniqueNoLock(t *testing.T) {
	var (
		c        = context.TODO()
		uniqueno = "1"
	)
	convey.Convey("AddUniqueNoLock", t, func(ctx convey.C) {
		succeed := d.AddUniqueNoLock(c, uniqueno)
		ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
			ctx.So(succeed, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddGrantUniqueLock(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
		mid   = int64(0)
	)
	convey.Convey("AddGrantUniqueLock", t, func(ctx convey.C) {
		succeed := d.AddGrantUniqueLock(c, token, mid)
		ctx.Convey("Then succeed should not be nil.", func(ctx convey.C) {
			ctx.So(succeed, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBranchCurrentCountCache(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
	)
	convey.Convey("BranchCurrentCountCache", t, func(ctx convey.C) {
		count, err := d.BranchCurrentCountCache(c, token)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetBranchCurrentCountCache(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
		count = int(0)
	)
	convey.Convey("SetBranchCurrentCountCache", t, func(ctx convey.C) {
		err := d.SetBranchCurrentCountCache(c, token, count)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoIncreaseBranchCurrentCountCache(t *testing.T) {
	var (
		c     = context.TODO()
		token = "1"
		count = uint64(0)
	)
	convey.Convey("IncreaseBranchCurrentCountCache", t, func(ctx convey.C) {
		err := d.IncreaseBranchCurrentCountCache(c, token, count)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCouponAllowanceCache(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		coupons = []*model.CouponAllowanceInfo{}
	)
	convey.Convey("CouponAllowanceCache", t, func(ctx convey.C) {
		err := d.SetCouponAllowanceCache(c, mid, 0, coupons)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		coupons, err := d.CouponAllowanceCache(c, mid, 0)
		ctx.Convey("Then err should be nil.coupons should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(coupons, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetPrizeCardCache(t *testing.T) {
	var (
		c               = context.TODO()
		mid       int64 = 1
		actID     int64 = 1
		prizeCard       = &model.PrizeCardRep{}
	)
	convey.Convey("SetPrizeCardCache ", t, func(ctx convey.C) {
		err := d.SetPrizeCardCache(c, mid, actID, prizeCard)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
func TestDaoSetPrizeCardsCache(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		actID      = int64(1)
		prizeCards = make([]*model.PrizeCardRep, 0)
		prizeCard  = &model.PrizeCardRep{}
	)
	prizeCards = append(prizeCards, prizeCard)
	convey.Convey("SetPrizeCardsCache ", t, func(ctx convey.C) {
		err := d.SetPrizeCardsCache(c, mid, actID, prizeCards)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPrizeCardCache(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(1)
		actID     = int64(1)
		ct        = int8(0)
		prizeCard = &model.PrizeCardRep{}
	)
	convey.Convey("PrizeCardCache", t, func(ctx convey.C) {
		d.SetPrizeCardCache(c, mid, actID, prizeCard)
		res, err := d.PrizeCardCache(c, mid, actID, ct)
		ctx.Convey("Then err should be nil.res should be not nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		d.DelPrizeCardKey(c, mid, actID, ct)
		res, err = d.PrizeCardCache(c, mid, actID, ct)
		ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
func TestDaoPrizeCardsCache(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		actID      = int64(1)
		prizeCards = make([]*model.PrizeCardRep, 0)
		prizeCard  = &model.PrizeCardRep{}
	)
	prizeCards = append(prizeCards, prizeCard)
	convey.Convey("PrizeCardCache", t, func(ctx convey.C) {
		d.SetPrizeCardsCache(c, mid, actID, prizeCards)
		res, err := d.PrizeCardsCache(c, mid, actID)
		ctx.Convey("Then err should be nil.res should be not nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		d.DelPrizeCardsKey(c, mid, actID)
		res, err = d.PrizeCardsCache(c, mid, actID)
		ctx.Convey("Then err should be nil.res should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldBeNil)
		})
	})
}
