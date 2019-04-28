package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/vip/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelBcoinSalary(t *testing.T) {
	convey.Convey("DelBcoinSalary", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			mid   = int64(1)
			month = xtime.Time(time.Now().Unix())
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelBcoinSalary(tx, mid, month)
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
	})
}

func TestDaoSelVipUserInfo(t *testing.T) {
	convey.Convey("SelVipUserInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			r, err := d.SelVipUserInfo(c, mid)
			ctx.Convey("Then err should be nil.r should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateVipUserInfo(t *testing.T) {
	convey.Convey("UpdateVipUserInfo", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			r     = &model.VipUserInfo{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.UpdateVipUserInfo(tx, r)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertVipChangeHistory(t *testing.T) {
	convey.Convey("InsertVipChangeHistory", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			r     = &model.VipChangeHistory{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			id, err := d.InsertVipChangeHistory(tx, r)
			if err == nil {
				if err = tx.Commit(); err != nil {
					tx.Rollback()
				}
			} else {
				tx.Rollback()
			}
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
