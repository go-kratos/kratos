package dao

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohitInfo(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("hitInfo", t, func(ctx convey.C) {
		p1 := hitInfo(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaohitChangeLog(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("hitChangeLog", t, func(ctx convey.C) {
		p1 := hitChangeLog(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaohitUser(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("hitUser", t, func(ctx convey.C) {
		p1 := hitUser(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaohitUserLog(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("hitUserLog", t, func(ctx convey.C) {
		p1 := hitUserLog(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBeginTran(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		p1, err := d.BeginTran(c)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponList(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		ct    = int8(0)
		no    = int64(0)
	)
	convey.Convey("CouponList", t, func(ctx convey.C) {
		_, err := d.CouponList(c, mid, state, ct, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCouponNoStartCheckList(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		ct    = int8(0)
		no    = int64(0)
	)
	convey.Convey("CouponNoStartCheckList", t, func(ctx convey.C) {
		_, err := d.CouponNoStartCheckList(c, mid, state, ct, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoBlanceNoStartCheckList(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ct  = int8(0)
		no  = int64(0)
	)
	convey.Convey("BlanceNoStartCheckList", t, func(ctx convey.C) {
		_, err := d.BlanceNoStartCheckList(c, mid, ct, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoByOrderNO(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(1)
		orderNO = "1235378892"
		ct      = int8(1)
	)
	convey.Convey("ByOrderNO", t, func(ctx convey.C) {
		r, err := d.ByOrderNO(c, mid, orderNO, ct)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateCouponInUse(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		cp  = &model.CouponInfo{}
		a   int64
		err error
	)
	convey.Convey("UpdateCouponInUse", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.UpdateCouponInUse(c, tx, cp)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoInsertPointHistory(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		l   = &model.CouponChangeLog{}
		err error
	)
	convey.Convey("InsertPointHistory", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err := d.InsertPointHistory(c, tx, l)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponInfo(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(27515401)
		token = "581807988720180417190545"
	)
	convey.Convey("CouponInfo", t, func(ctx convey.C) {
		r, err := d.CouponInfo(c, mid, token)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountByState(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		no    = int64(0)
		stime = time.Now()
	)
	convey.Convey("CountByState", t, func(ctx convey.C) {
		count, err := d.CountByState(c, mid, state, no, stime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponPage(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		no    = int64(0)
		start = int(0)
		ps    = int(0)
		stime = time.Now()
	)
	convey.Convey("CouponPage", t, func(ctx convey.C) {
		_, err := d.CouponPage(c, mid, state, no, start, ps, stime)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddCoupon(t *testing.T) {
	var (
		c  = context.TODO()
		cp = &model.CouponInfo{
			Mid:         1,
			CouponToken: token(),
			State:       0,
			StartTime:   time.Now().Unix(),
			ExpireTime:  time.Now().AddDate(0, 0, 1).Unix(),
		}
	)
	convey.Convey("AddCoupon", t, func(ctx convey.C) {
		a, err := d.AddCoupon(c, cp)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchAddCoupon(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		mid = int64(0)
		cps = []*model.CouponInfo{}
		err error
		a   int64
	)
	convey.Convey("BatchAddCoupon", t, func(ctx convey.C) {
		cp := &model.CouponInfo{}
		cp.CouponToken = token()
		cp.Mid = mid
		cp.State = model.NotUsed
		cp.StartTime = time.Now().Unix()
		cp.ExpireTime = time.Now().AddDate(0, 0, 2).Unix()
		cp.Origin = 1
		cp.CouponType = 1
		cp.CTime = xtime.Time(time.Now().Unix())
		cp.BatchToken = "1234"
		cps = append(cps, cp)
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.BatchAddCoupon(c, tx, mid, cps)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateCoupon(t *testing.T) {
	var (
		c           = context.TODO()
		mid         = int64(27515800)
		state       = int8(0)
		useVer      = int64(0)
		ver         = int64(1)
		couponToken = "510204683920180420110002"
	)
	convey.Convey("UpdateCoupon", t, func(ctx convey.C) {
		_, err := d.UpdateCoupon(c, mid, state, useVer, ver, couponToken)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoByThirdTradeNo(t *testing.T) {
	var (
		c            = context.TODO()
		thirdTradeNo = "12156121892"
		ct           = int8(2)
	)
	convey.Convey("ByThirdTradeNo", t, func(ctx convey.C) {
		r, err := d.ByThirdTradeNo(c, thirdTradeNo, ct)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponBlances(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(520)
		ct  = int8(2)
		no  = int64(0)
	)
	convey.Convey("CouponBlances", t, func(ctx convey.C) {
		_, err := d.CouponBlances(c, mid, ct, no)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateBlance(t *testing.T) {
	var (
		c       = context.TODO()
		tx      = &sql.Tx{}
		id      = int64(0)
		mid     = int64(0)
		ver     = int64(0)
		balance = int64(0)
		a       int64
		err     error
	)
	convey.Convey("UpdateBlance", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.UpdateBlance(c, tx, id, mid, ver, balance)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchUpdateBlance(t *testing.T) {
	var (
		c       = context.TODO()
		tx      = &sql.Tx{}
		mid     = int64(1)
		blances = []*model.CouponBalanceInfo{}
		err     error
	)
	blances = append(blances, &model.CouponBalanceInfo{
		ID:      116197,
		Balance: 1,
		Ver:     0,
	})
	convey.Convey("BatchUpdateBlance", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		_, err = d.BatchUpdateBlance(c, tx, mid, blances)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		err = tx.Commit()
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestDaoBatchInsertBlanceLog(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		mid = int64(0)
		ls  = []*model.CouponBalanceChangeLog{}
		err error
	)
	convey.Convey("BatchInsertBlanceLog", t, func(ctx convey.C) {
		blog := new(model.CouponBalanceChangeLog)
		blog.OrderNo = "11"
		blog.Mid = mid
		blog.BatchToken = "123"
		blog.ChangeType = model.Consume
		blog.Ctime = xtime.Time(time.Now().Unix())
		ls = append(ls, blog)
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		_, err = d.BatchInsertBlanceLog(c, tx, mid, ls)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		err = tx.Commit()
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestDaoAddOrder(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		o   = &model.CouponOrder{}
		a   int64
		err error
	)
	convey.Convey("AddOrder", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.AddOrder(c, tx, o)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddOrderLog(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		o   = &model.CouponOrderLog{}
		a   int64
		err error
	)
	convey.Convey("AddOrderLog", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.AddOrderLog(c, tx, o)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponCarToonCount(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		no    = int64(0)
		ct    = int8(0)
		state = int8(0)
		stime = time.Now()
	)
	convey.Convey("CouponCarToonCount", t, func(ctx convey.C) {
		count, err := d.CouponCarToonCount(c, mid, no, ct, state, stime)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCouponNotUsedPage(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		ct    = int8(0)
		no    = int64(0)
		stime = time.Now()
		pn    = int(0)
		ps    = int(0)
	)
	convey.Convey("CouponNotUsedPage", t, func(ctx convey.C) {
		_, err := d.CouponNotUsedPage(c, mid, ct, no, stime, pn, ps)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCouponExpirePage(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(27515301)
		ct    = int8(1)
		no    = time.Now().Unix()
		stime = time.Now().AddDate(-1, 0, 0)
		pn    = int(1)
		ps    = int(10)
	)
	convey.Convey("CouponExpirePage", t, func(ctx convey.C) {
		_, err := d.CouponExpirePage(c, mid, ct, no, stime, pn, ps)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoOrderUsedPage(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int8(0)
		ct    = int8(0)
		stime = time.Now()
		pn    = int(0)
		ps    = int(0)
	)
	convey.Convey("OrderUsedPage", t, func(ctx convey.C) {
		_, err := d.OrderUsedPage(c, mid, state, ct, stime, pn, ps)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddBalanceCoupon(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		b   = &model.CouponBalanceInfo{}
		a   int64
		err error
	)
	convey.Convey("AddBalanceCoupon", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.AddBalanceCoupon(c, tx, b)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoByMidAndBatchToken(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(1)
		batchToken = "441539420220180806174505"
	)
	convey.Convey("ByMidAndBatchToken", t, func(ctx convey.C) {
		_, err := d.ByMidAndBatchToken(c, mid, batchToken)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddBalanceChangeLog(t *testing.T) {
	var (
		c   = context.TODO()
		tx  = &sql.Tx{}
		bl  = &model.CouponBalanceChangeLog{}
		a   int64
		err error
	)
	convey.Convey("AddBalanceChangeLog", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.AddBalanceChangeLog(c, tx, bl)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchInfo(t *testing.T) {
	var (
		c     = context.TODO()
		token = "900364604420180912170927"
	)
	convey.Convey("BatchInfo", t, func(ctx convey.C) {
		r, err := d.BatchInfo(c, token)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateBatchInfo(t *testing.T) {
	var (
		c     = context.TODO()
		tx    = &sql.Tx{}
		token = ""
		count = int(0)
		a     int64
		err   error
	)
	convey.Convey("UpdateBatchInfo", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.UpdateBatchInfo(c, tx, token, count)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateBatchLimitInfo(t *testing.T) {
	var (
		c     = context.TODO()
		tx    = &sql.Tx{}
		token = ""
		count = int(0)
		a     int64
		err   error
	)
	convey.Convey("UpdateBatchLimitInfo", t, func(ctx convey.C) {
		tx, err = d.BeginTran(c)
		ctx.So(err, convey.ShouldBeNil)
		a, err = d.UpdateBatchLimitInfo(c, tx, token, count)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGrantCouponLog(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		token = ""
		ct    = int8(0)
	)
	convey.Convey("GrantCouponLog", t, func(ctx convey.C) {
		_, err := d.GrantCouponLog(c, mid, token, ct)
		ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAllBranchInfo(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("AllBranchInfo", t, func(ctx convey.C) {
		res, err := d.AllBranchInfo(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoCountByBranchToken(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		token = ""
	)
	convey.Convey("CountByBranchToken", t, func(ctx convey.C) {
		count, err := d.CountByBranchToken(c, mid, token)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

// get coupon token
func token() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("20060102150405"))
	return b.String()
}

func TestDaoReceiveLog(t *testing.T) {
	var (
		c            = context.Background()
		appkey       = "7c7ac0db1aa05587"
		orderNo      = "1536657724"
		ct      int8 = 3
	)
	convey.Convey("ReceiveLog ", t, func(ctx convey.C) {
		r, err := d.ReceiveLog(c, appkey, orderNo, ct)
		ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		})
		r, err = d.ReceiveLog(c, "", "", 21)
		ctx.Convey("Then err should be nil.r should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxAddReceiveLog(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		rlog  = &model.CouponReceiveLog{Appkey: fmt.Sprintf("%d", time.Now().Unix()), CouponType: int8(rand.Int63n(127))}
	)
	convey.Convey("TxAddReceiveLog ", t, func(ctx convey.C) {
		err := d.TxAddReceiveLog(tx, rlog)
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
