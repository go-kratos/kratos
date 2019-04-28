package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeUpAccountCount(t *testing.T) {
	convey.Convey("UpAccountCount", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			query     = ""
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.UpAccountCount(c, query, isDeleted)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeListUpAccount(t *testing.T) {
	convey.Convey("ListUpAccount", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			query     = ""
			isDeleted = int(0)
			from      = int(0)
			limit     = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.ListUpAccount(c, query, isDeleted, from, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetUpAccount(t *testing.T) {
	convey.Convey("GetUpAccount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_account(mid) VALUES(1001)")
			up, err := d.GetUpAccount(c, mid)
			ctx.Convey("Then err should be nil.up should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(up, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeTxBreachUpAccount(t *testing.T) {
	convey.Convey("TxBreachUpAccount", t, func(ctx convey.C) {
		var (
			tx, _      = d.BeginTran(context.Background())
			total      = int64(0)
			unwithdraw = int64(0)
			mid        = int64(1001)
			newVersion = int64(0)
			oldVersion = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxBreachUpAccount(tx, total, unwithdraw, mid, newVersion, oldVersion)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
