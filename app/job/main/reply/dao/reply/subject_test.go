package reply

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyUpMeta(t *testing.T) {
	convey.Convey("UpMeta", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			tp   = int8(0)
			meta = ""
			now  = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.UpMeta(c, oid, tp, meta, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyTxUpMeta(t *testing.T) {
	convey.Convey("TxUpMeta", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			meta  = ""
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxUpMeta(tx, oid, tp, meta, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxUpAttr(t *testing.T) {
	convey.Convey("TxUpAttr", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			attr  = uint32(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxUpAttr(tx, oid, tp, attr, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrCount(t *testing.T) {
	convey.Convey("TxIncrCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxIncrCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrFCount(t *testing.T) {
	convey.Convey("TxIncrFCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxIncrFCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrMCount(t *testing.T) {
	convey.Convey("TxIncrMCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxIncrMCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxDecrMCount(t *testing.T) {
	convey.Convey("TxDecrMCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxDecrMCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrRCount(t *testing.T) {
	convey.Convey("TxIncrRCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxIncrRCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxDecrCount(t *testing.T) {
	convey.Convey("TxDecrCount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxDecrCount(tx, oid, tp, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxIncrACount(t *testing.T) {
	convey.Convey("TxIncrACount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			count = int(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxIncrACount(tx, oid, tp, count, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyTxDecrACount(t *testing.T) {
	convey.Convey("TxDecrACount", t, func(ctx convey.C) {
		var (
			tx, _ = d.mysql.Begin(context.Background())
			oid   = int64(0)
			tp    = int8(0)
			count = int(0)
			now   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.Subject.TxDecrACount(tx, oid, tp, count, now)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
		tx.Rollback()
	})
}

func TestReplyGet2(t *testing.T) {
	convey.Convey("Get", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			sub, err := d.Subject.Get(c, oid, tp)
			ctx.Convey("Then err should be nil.sub should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sub, convey.ShouldNotBeNil)
			})
		})
	})
}
