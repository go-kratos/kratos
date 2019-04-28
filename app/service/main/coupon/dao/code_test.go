package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/coupon/model"
	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCouponCode(t *testing.T) {
	convey.Convey("CouponCode", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			code = "2asazxcvfdsb"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CouponCode(c, code)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCountCodeByMid(t *testing.T) {
	convey.Convey("CountCodeByMid", t, func(convCtx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(1)
			batckToken = "allowance_batch100"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			count, err := d.CountCodeByMid(c, mid, batckToken)
			convCtx.Convey("Then err should be nil.count should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateCodeState(t *testing.T) {
	convey.Convey("TxUpdateCodeState", t, func(convCtx convey.C) {
		var (
			tx = &sql.Tx{}
			a  = &model.CouponCode{
				State: 2,
				Mid:   1,
				Ver:   1,
			}
			err error
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			tx, err = d.BeginTran(context.Background())
			convCtx.So(err, convey.ShouldBeNil)
			aff, err := d.TxUpdateCodeState(tx, a)
			convCtx.Convey("Then err should be nil.aff should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(aff, convey.ShouldNotBeNil)
			})
		})
	})
}
