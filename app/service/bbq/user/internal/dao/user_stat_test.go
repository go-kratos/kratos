package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawBatchUserStatistics(t *testing.T) {
	convey.Convey("RawBatchUserStatistics", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{88895104}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawBatchUserStatistics(c, mids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxIncrUserStatisticsFollow(t *testing.T) {
	convey.Convey("TxIncrUserStatisticsFollow", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxIncrUserStatisticsFollow(tx, mid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxIncrUserStatisticsFan(t *testing.T) {
	convey.Convey("TxIncrUserStatisticsFan", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxIncrUserStatisticsFan(tx, mid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDecrUserStatisticsFollow(t *testing.T) {
	convey.Convey("TxDecrUserStatisticsFollow", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxDecrUserStatisticsFollow(tx, mid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDecrUserStatisticsFan(t *testing.T) {
	convey.Convey("TxDecrUserStatisticsFan", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(88895104)
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num, err := d.TxDecrUserStatisticsFan(tx, mid)
			ctx.Convey("Then err should be nil.num should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(num, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxIncrUserStatisticsField(t *testing.T) {
	convey.Convey("TxIncrUserStatisticsField", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895104)
			field = "like_total"
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rowsAffected, err := d.TxIncrUserStatisticsField(c, tx, mid, field)
			ctx.Convey("Then err should be nil.rowsAffected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rowsAffected, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxDescUserStatisticsField(t *testing.T) {
	convey.Convey("TxDescUserStatisticsField", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(88895104)
			field = "like_total"
		)
		tx, _ := d.BeginTran(c)
		defer func() {
			tx.Rollback()
		}()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rowsAffected, err := d.TxDescUserStatisticsField(c, tx, mid, field)
			ctx.Convey("Then err should be nil.rowsAffected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rowsAffected, convey.ShouldNotBeNil)
			})
		})
	})
}
