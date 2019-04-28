package dao

import (
	"context"
	"go-common/app/job/main/vip/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSalaryVideoCouponList(t *testing.T) {
	convey.Convey("SalaryVideoCouponList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dv  = "2018_09"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.SalaryVideoCouponList(c, mid, dv)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddSalaryLog(t *testing.T) {
	convey.Convey("AddSalaryLog", t, func(ctx convey.C) {
		var (
			c = context.Background()
			l = &model.VideoCouponSalaryLog{
				Mid: time.Now().Unix(),
			}
			dv = "2018_09"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddSalaryLog(c, l, dv)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
