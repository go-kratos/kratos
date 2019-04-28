package dao

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"time"
	// xsql"database/sql"
	"go-common/app/admin/main/coupon/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

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

func TestDaohitAllowanceInfo(t *testing.T) {
	convey.Convey("hitAllowanceInfo", t, func(convCtx convey.C) {
		var (
			mid = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitAllowanceInfo(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitAllowanceChangeLog(t *testing.T) {
	convey.Convey("hitAllowanceChangeLog", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitAllowanceChangeLog(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaohitViewInfo(t *testing.T) {
	convey.Convey("hitViewInfo", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := hitViewInfo(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

// go test  -test.v -test.run TestBatchList
func TestDaoBatchList(t *testing.T) {
	convey.Convey("BatchList", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			appid = int64(1)
			t     = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BatchList(c, appid, t)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchViewList(t *testing.T) {
	convey.Convey("BatchViewList", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			appid      = int64(0)
			batchToken = ""
			no         = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.BatchViewList(c, appid, batchToken, no)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBatchInfo(t *testing.T) {
	convey.Convey("AddBatchInfo", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			bi = &model.CouponBatchInfo{}
			b  bytes.Buffer
		)
		b.WriteString(fmt.Sprintf("%07d", rand.Int63n(9999999)))
		b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
		b.WriteString(time.Now().Format("20060102150405"))
		bi.BatchToken = b.String()
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.AddBatchInfo(c, bi)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllAppInfo(t *testing.T) {
	convey.Convey("AllAppInfo", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.AllAppInfo(c)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAllowanceBatchInfo(t *testing.T) {
	convey.Convey("AddAllowanceBatchInfo", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			b  bytes.Buffer
			bi = &model.CouponBatchInfo{}
		)
		b.WriteString(fmt.Sprintf("%07d", rand.Int63n(9999999)))
		b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
		b.WriteString(time.Now().Format("20060102150405"))
		bi.BatchToken = b.String()
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.AddAllowanceBatchInfo(c, bi)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAllowanceBatchInfo(t *testing.T) {
	convey.Convey("UpdateAllowanceBatchInfo", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			b = &model.CouponBatchInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateAllowanceBatchInfo(c, b)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBatchStatus(t *testing.T) {
	convey.Convey("UpdateBatchStatus", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			status   = int8(0)
			operator = ""
			id       = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateBatchStatus(c, status, operator, id)
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchInfo(t *testing.T) {
	convey.Convey("BatchInfo", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			token = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.BatchInfo(c, token)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateAllowanceStatus(t *testing.T) {
	convey.Convey("UpdateAllowanceStatus", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			state = int8(0)
			mid   = int64(0)
			token = ""
			ver   = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateAllowanceStatus(c, tx, state, mid, token, ver)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllowanceByToken(t *testing.T) {
	convey.Convey("AllowanceByToken", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(13)
			token = "000000119720180929180009"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.AllowanceByToken(c, mid, token)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertCouponAllowanceHistory(t *testing.T) {
	convey.Convey("InsertCouponAllowanceHistory", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			l     = &model.CouponAllowanceChangeLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.InsertCouponAllowanceHistory(c, tx, l)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAllowanceList(t *testing.T) {
	convey.Convey("AllowanceList", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgAllowanceSearch{Mid: 3}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.AllowanceList(c, arg)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddViewBatch(t *testing.T) {
	convey.Convey("AddViewBatch", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgCouponViewBatch{}
			b   bytes.Buffer
		)
		b.WriteString(fmt.Sprintf("%07d", rand.Int63n(9999999)))
		b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
		b.WriteString(time.Now().Format("20060102150405"))
		arg.BatchToken = b.String()
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddViewBatch(c, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateViewBatch(t *testing.T) {
	convey.Convey("UpdateViewBatch", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgCouponViewBatch{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateViewBatch(c, arg)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxUpdateViewInfo(t *testing.T) {
	convey.Convey("TxUpdateViewInfo", t, func(convCtx convey.C) {
		var (
			tx, _       = d.BeginTran(context.Background())
			status      = int8(0)
			couponToken = ""
			mid         = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.TxUpdateViewInfo(tx, status, couponToken, mid)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoTxCouponViewLog(t *testing.T) {
	convey.Convey("TxCouponViewLog", t, func(convCtx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			arg   = &model.CouponChangeLog{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.TxCouponViewLog(tx, arg)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCouponViewInfo(t *testing.T) {
	convey.Convey("CouponViewInfo", t, func(convCtx convey.C) {
		var (
			c           = context.Background()
			couponToken = ""
			mid         = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.CouponViewInfo(c, couponToken, mid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSearchViewCouponCount(t *testing.T) {
	convey.Convey("SearchViewCouponCount", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgSearchCouponView{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.SearchViewCouponCount(c, arg)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSearchViewCouponInfo(t *testing.T) {
	convey.Convey("SearchViewCouponInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			arg = &model.ArgSearchCouponView{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.SearchViewCouponInfo(c, arg)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBatchAddAllowanceCoupon(t *testing.T) {
	convey.Convey("BatchAddAllowanceCoupon", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			b     bytes.Buffer
			tx, _ = d.BeginTran(context.Background())
			cps   = []*model.CouponAllowanceInfo{}
		)
		b.WriteString(fmt.Sprintf("%05d", 1))
		b.WriteString(fmt.Sprintf("%02d", rand.Int63n(99)))
		b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
		b.WriteString(time.Now().Format("20060102150405"))
		cp := &model.CouponAllowanceInfo{CouponToken: b.String()}
		cps = append(cps, cp)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.BatchAddAllowanceCoupon(c, tx, cps)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBatchInfo(t *testing.T) {
	convey.Convey("UpdateBatchInfo", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			tx, _ = d.BeginTran(context.Background())
			token = ""
			count = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			a, err := d.UpdateBatchInfo(c, tx, token, count)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			convCtx.Convey("Then err should be nil.a should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}
