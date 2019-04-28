package dao

import (
	"context"
	"go-common/app/admin/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertCreditRecord(t *testing.T) {
	convey.Convey("InsertCreditRecord", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			cr = &model.CreditRecord{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InsertCreditRecord(c, cr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertCreditRecord(t *testing.T) {
	convey.Convey("TxInsertCreditRecord", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			cr    = &model.CreditRecord{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxInsertCreditRecord(tx, cr)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreditRecords(t *testing.T) {
	convey.Convey("CreditRecords", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			crs, err := d.CreditRecords(c, mid)
			ctx.Convey("Then err should be nil.crs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(crs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeductedScore(t *testing.T) {
	convey.Convey("DeductedScore", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			deducted, err := d.DeductedScore(c, id)
			ctx.Convey("Then err should be nil.deducted should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(deducted, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertCreditScore(t *testing.T) {
	convey.Convey("InsertCreditScore", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(100)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "DELETE FROM credit_score WHERE mid = 100")
			rows, err := d.InsertCreditScore(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateCreditScore(t *testing.T) {
	convey.Convey("UpdateCreditScore", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			score = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateCreditScore(c, mid, score)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxUpdateCreditScore(t *testing.T) {
	convey.Convey("TxUpdateCreditScore", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			mid   = int64(0)
			score = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxUpdateCreditScore(tx, mid, score)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreditScore(t *testing.T) {
	convey.Convey("CreditScore", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			score, err := d.CreditScore(c, mid)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreditScores(t *testing.T) {
	convey.Convey("CreditScores", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			scores, err := d.CreditScores(c, mids)
			ctx.Convey("Then err should be nil.scores should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(scores, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxRecoverCreditScore(t *testing.T) {
	convey.Convey("TxRecoverCreditScore", t, func(ctx convey.C) {
		var (
			tx, _    = d.BeginTran(context.Background())
			deducted = int(0)
			mid      = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			rows, err := d.TxRecoverCreditScore(tx, deducted, mid)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
