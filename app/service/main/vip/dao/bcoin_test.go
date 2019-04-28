package dao

import (
	"context"
	"go-common/app/service/main/vip/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBcoinSalaryList(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		start  = time.Now()
		end    = time.Now()
		expRes []*model.VipBcoinSalary
	)
	convey.Convey("BcoinSalaryList", t, func(ctx convey.C) {

		res, err := d.BcoinSalaryList(c, mid, start, end)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldEqual, expRes)
		})
	})
}

func TestDaoSelLastBcoin(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(20606508)
	)
	convey.Convey("SelLastBcoin", t, func(ctx convey.C) {
		_, err := d.SelLastBcoin(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoInsertVipBcoinSalary(t *testing.T) {
	var (
		c = context.TODO()
		r = &model.VipBcoinSalary{}
	)
	convey.Convey("InsertVipBcoinSalary", t, func(ctx convey.C) {
		err := d.InsertVipBcoinSalary(c, r)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
