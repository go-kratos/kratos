package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagDelArchiveRatio(t *testing.T) {
	convey.Convey("DelArchiveRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = int(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelArchiveRatio(c, ctype, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagDelUpRatio(t *testing.T) {
	convey.Convey("DelUpRatio", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ctype = int(0)
			limit = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelUpRatio(c, ctype, limit)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagInsertAvRatio(t *testing.T) {
	convey.Convey("InsertAvRatio", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "DELETE FROM av_charge_ratio WHERE av_id=2")
			rows, err := d.InsertAvRatio(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagInsertUpRatio(t *testing.T) {
	convey.Convey("InsertUpRatio", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "DELETE FROM up_charge_ratio")
			rows, err := d.InsertUpRatio(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
