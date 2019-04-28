package dao

import (
	"context"
	// "database/sql"
	"go-common/app/job/main/coupon/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitInfo(t *testing.T) {
	convey.Convey("hitInfo", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitInfo(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitChangeLog(t *testing.T) {
	convey.Convey("hitChangeLog", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitChangeLog(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitUser(t *testing.T) {
	convey.Convey("hitUser", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitUser(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitUserLog(t *testing.T) {
	convey.Convey("hitUserLog", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitUserLog(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

// go test  -test.v -test.run TestDaoUpdateCoupon
func TestDaoUpdateCoupon(t *testing.T) {
	convey.Convey("UpdateCoupon", t, func(convCtx convey.C) {
		var (
			c           = context.Background()
			tx, _       = d.BeginTran(context.Background())
			mid         = int64(1)
			state       = int8(1)
			useVer      = int64(11)
			ver         = int64(2)
			couponToken = "729792667120180402161647"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateCoupon(c, tx, mid, state, useVer, ver, couponToken)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoCouponInfo(t *testing.T) {
	convey.Convey("CouponInfo", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			token = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.CouponInfo(c, mid, token)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCouponList(t *testing.T) {
	convey.Convey("CouponList", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			index = int64(0)
			state = int8(0)
			no, _ = time.Parse("2006-01-02 15:04:05", "2018-12-27 17:28:51")
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CouponList(c, index, state, no)
			convCtx.Convey("Then err should be nil.res should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoInsertPointHistory(t *testing.T) {
	convey.Convey("InsertPointHistory", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			l     = &model.CouponChangeLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.InsertPointHistory(c, tx, l)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := d.BeginTran(c)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoByOrderNo(t *testing.T) {
	convey.Convey("ByOrderNo", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			orderNo = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.ByOrderNo(c, orderNo)
			convCtx.Convey("Then err should be nil.r should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateOrderState(t *testing.T) {
	convey.Convey("UpdateOrderState", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(context.Background())
			mid     = int64(0)
			state   = int8(0)
			useVer  = int64(0)
			ver     = int64(0)
			orderNo = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateOrderState(c, tx, mid, state, useVer, ver, orderNo)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoAddOrderLog(t *testing.T) {
	convey.Convey("AddOrderLog", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			o     = &model.CouponOrderLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.AddOrderLog(c, tx, o)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoConsumeCouponLog(t *testing.T) {
	convey.Convey("ConsumeCouponLog", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			orderNo = ""
			ct      = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rs, err := d.ConsumeCouponLog(c, mid, orderNo, ct)
			convCtx.Convey("Then err should be nil.rs should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoByMidAndBatchToken(t *testing.T) {
	convey.Convey("ByMidAndBatchToken", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(0)
			batchToken = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.ByMidAndBatchToken(c, mid, batchToken)
			convCtx.Convey("Then err should be nil.r should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateBlance(t *testing.T) {
	convey.Convey("UpdateBlance", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(context.Background())
			id      = int64(0)
			mid     = int64(0)
			ver     = int64(0)
			balance = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateBlance(c, tx, id, mid, ver, balance)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoOrderInPay(t *testing.T) {
	convey.Convey("OrderInPay", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			state = int8(0)
			no    = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.OrderInPay(c, state, no)
			convCtx.Convey("Then err should be nil.res should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBatchUpdateBlance(t *testing.T) {
	convey.Convey("BatchUpdateBlance", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			tx, _   = d.BeginTran(context.Background())
			mid     = int64(0)
			blances = []*model.CouponBalanceInfo{}
			blance  = &model.CouponBalanceInfo{}
		)
		blances = append(blances, blance)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.BatchUpdateBlance(c, tx, mid, blances)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoBatchInsertBlanceLog(t *testing.T) {
	convey.Convey("BatchInsertBlanceLog", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			mid   = int64(0)
			ls    = []*model.CouponBalanceChangeLog{}
			l     = &model.CouponBalanceChangeLog{}
		)
		ls = append(ls, l)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.BatchInsertBlanceLog(c, tx, mid, ls)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoBlanceList(t *testing.T) {
	convey.Convey("BlanceList", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ct  = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BlanceList(c, mid, ct)
			convCtx.Convey("Then err should be nil.res should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateUserCard(t *testing.T) {
	convey.Convey("UpdateUserCard", t, func(convCtx convey.C) {
		var (
			c           = context.Background()
			mid         = int64(0)
			state       = int8(0)
			couponToken = ""
			batchToken  = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateUserCard(c, mid, state, couponToken, batchToken)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoClose(t *testing.T) {
	convey.Convey("TestDaoClose", t, func(convCtx convey.C) {
		d.Close()
	})
}
