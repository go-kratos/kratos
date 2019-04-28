package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/library/time"
)

func TestDao_TxUpDelay(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		dt    time.Time
	)
	Convey("TxUpArchiveState", t, func(ctx C) {
		_, err := d.TxUpDelay(tx, 123, 23333, 0, 0, dt)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxDelDelay(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxDelDelay", t, func(ctx C) {
		_, err := d.TxDelDelay(tx, 123, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestArchiveDelay(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(23333)
		tp  = int8(0)
	)
	Convey("Delay", t, func(ctx C) {
		_, err := d.Delay(c, aid, tp)
		ctx.Convey("Then err should be nil.dl should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}
