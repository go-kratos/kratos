package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeBreachCount(t *testing.T) {
	convey.Convey("BreachCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.BreachCount(c, query)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeListArchiveBreach(t *testing.T) {
	convey.Convey("ListArchiveBreach", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "id > 0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			breachs, err := d.ListArchiveBreach(c, query)
			ctx.Convey("Then err should be nil.breachs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(breachs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeTxInsertAvBreach(t *testing.T) {
	convey.Convey("TxInsertAvBreach", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			val   = "(520,1100,'2018-01-01',10,0,'aa','2018-01-01')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(context.Background(), "DELETE FROM av_breach_record WHERE av_id = 520")
			rows, err := d.TxInsertAvBreach(tx, val)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetAvBreachByMIDs(t *testing.T) {
	convey.Convey("GetAvBreachByMIDs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mids  = []int64{1000}
			types = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_breach_record(mid, ctype) VALUES(1000, 1)")
			breachs, err := d.GetAvBreachByMIDs(c, mids, types)
			ctx.Convey("Then err should be nil.breachs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(breachs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeTxUpdateBreachPre(t *testing.T) {
	convey.Convey("TxUpdateBreachPre", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			aids  = []int64{1001}
			cdate = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxUpdateBreachPre(tx, aids, cdate)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
