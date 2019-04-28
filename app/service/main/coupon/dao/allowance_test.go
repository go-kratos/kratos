package dao

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitAllowanceInfo(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("hitAllowanceInfo", t, func(ctx convey.C) {
		p1 := hitAllowanceInfo(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaohitAllowanceChangeLog(t *testing.T) {
	var (
		mid = int64(1)
	)
	convey.Convey("hitAllowanceChangeLog", t, func(ctx convey.C) {
		p1 := hitAllowanceChangeLog(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoByStateAndExpireAllowances(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(1)
		state = int8(0)
		no    = int64(0)
	)
	convey.Convey("ByStateAndExpireAllowances", t, func(ctx convey.C) {
		res, err := d.ByStateAndExpireAllowances(c, mid, state, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAllowanceByOrderNO(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		orderNO = "1"
	)
	convey.Convey("AllowanceByOrderNO", t, func(ctx convey.C) {
		_, err := d.AllowanceByOrderNO(c, mid, orderNO)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUsableAllowances(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(1)
		state = int8(0)
		no    = int64(0)
	)
	convey.Convey("UsableAllowances", t, func(ctx convey.C) {
		_, err := d.UsableAllowances(c, mid, state, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAllowanceByToken(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		token = "1"
	)
	convey.Convey("AllowanceByToken", t, func(ctx convey.C) {
		_, err := d.AllowanceByToken(c, mid, token)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateAllowanceCouponInUse(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		cp  = &model.CouponAllowanceInfo{}
		err error
	)
	convey.Convey("UpdateAllowanceCouponInUse", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err := d.UpdateAllowanceCouponInUse(c, tx, cp)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateAllowanceCouponToUse(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		cp    = &model.CouponAllowanceInfo{}
	)
	convey.Convey("UpdateAllowanceCouponToUse ", t, func(ctx convey.C) {
		a, err := d.UpdateAllowanceCouponToUse(c, tx, cp)
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoUpdateAllowanceCouponToUsed(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		cp    = &model.CouponAllowanceInfo{}
	)
	convey.Convey("UpdateAllowanceCouponToUsed ", t, func(ctx convey.C) {
		a, err := d.UpdateAllowanceCouponToUsed(c, tx, cp)
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestDaoInsertCouponAllowanceHistory(t *testing.T) {
	var (
		c  = context.TODO()
		tx = &sql.Tx{}
		l  = &model.CouponAllowanceChangeLog{
			CouponToken: token(),
			OrderNO:     "1",
		}
		err error
	)
	convey.Convey("InsertCouponAllowanceHistory", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err := d.InsertCouponAllowanceHistory(c, tx, l)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountByAllowanceBranchToken(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		token = "allowance_test1"
	)
	convey.Convey("CountByAllowanceBranchToken", t, func(ctx convey.C) {
		count, err := d.CountByAllowanceBranchToken(c, mid, token)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGetCouponByOrderNo(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		orderNo = "allowance_test1"
	)
	convey.Convey("GetCouponByOrderNo ", t, func(ctx convey.C) {
		res, err := d.GetCouponByOrderNo(c, mid, orderNo)
		ctx.Convey("Then err should be not nil.res should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTxAddAllowanceCoupon(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		cp    = &model.CouponAllowanceInfo{CouponToken: strconv.FormatInt(rand.Int63n(999999), 10)}
	)
	convey.Convey("TxAddAllowanceCoupon ", t, func(ctx convey.C) {
		err := d.TxAddAllowanceCoupon(tx, cp)
		if err == nil {
			if err = tx.Commit(); err != nil {
				tx.Rollback()
			}
		} else {
			tx.Rollback()
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBatchAddAllowanceCoupon(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		mid = int64(0)
		cps = []*model.CouponAllowanceInfo{}
		err error
	)
	convey.Convey("BatchAddAllowanceCoupon", t, func(ctx convey.C) {
		cps = append(cps, &model.CouponAllowanceInfo{
			CouponToken: token(),
			Mid:         mid,
			State:       model.NotUsed,
			StartTime:   time.Now().Unix(),
			ExpireTime:  time.Now().AddDate(0, 0, 10).Unix(),
			Origin:      int64(1),
			CTime:       xtime.Time(time.Now().Unix()),
			BatchToken:  "1",
			Amount:      float64(1),
			FullAmount:  float64(100),
			AppID:       int64(1),
		})
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		_, err = d.BatchAddAllowanceCoupon(c, tx, mid, cps)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAllowanceList(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		no    = int64(0)
		stime = time.Now()
	)
	convey.Convey("AllowanceList", t, func(ctx convey.C) {
		_, err := d.AllowanceList(c, mid, state, no, stime)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
