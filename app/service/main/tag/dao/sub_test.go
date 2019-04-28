package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoshardSub(t *testing.T) {
	convey.Convey("shardSub", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.shardSub(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddSub(t *testing.T) {
	convey.Convey("TxAddSub", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
			tid = int64(1833)
		)
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := d.TxAddSub(tx, mid, tid)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxDelSub(t *testing.T) {
	convey.Convey("TxDelSub", t, func(ctx convey.C) {
		var (
			mid = int64(35152246)
			tid = int64(1833)
		)
		tx, err := d.BeginTran(context.Background())
		if err != nil {
			return
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.TxDelSub(tx, mid, tid)
		})
		tx.Commit()
	})
}

func TestDaoSub(t *testing.T) {
	convey.Convey("Sub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, rem, err := d.Sub(c, mid)
			ctx.Convey("Then err should be nil.res,rem should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rem, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSubList(t *testing.T) {
	convey.Convey("SubList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(35152246)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SubList(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
