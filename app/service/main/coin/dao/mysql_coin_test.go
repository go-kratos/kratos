package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateCoin(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(0)
		coin = float64(0)
	)
	convey.Convey("UpdateCoin", t, func(ctx convey.C) {
		err := d.UpdateCoin(c, mid, coin)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoTxUpdateCoins(t *testing.T) {
	var (
		c     = context.TODO()
		tx, _ = d.BeginTran(c)
		mid   = int64(1)
		coin  = float64(20)
	)
	convey.Convey("TxUpdateCoins", t, func(ctx convey.C) {
		err := d.TxUpdateCoins(c, tx, mid, coin)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			if err != nil {
				tx.Rollback()
			}
			ctx.So(err, convey.ShouldBeNil)
			tx.Commit()
		})
	})
}

func TestDaoTxUserCoin(t *testing.T) {
	var (
		c     = context.TODO()
		tx, _ = d.BeginTran(c)
		mid   = int64(0)
	)
	convey.Convey("TxUserCoin", t, func(ctx convey.C) {
		count, err := d.TxUserCoin(c, tx, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			if err != nil {
				tx.Rollback()
			}
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("count should not be nil", func(ctx convey.C) {
			if count != 0 {
				tx.Commit()
			}
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRawUserCoin(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("RawUserCoin", t, func(ctx convey.C) {
		res, err := d.RawUserCoin(c, mid)
		ctx.Convey("Error should be nil", func(ctx convey.C) {

			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("res should not be nil", func(ctx convey.C) {
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
