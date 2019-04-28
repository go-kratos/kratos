package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertUpBillBatch(t *testing.T) {
	convey.Convey("InsertUpBillBatch", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			values = "(1,2,3,4,5,6,7,8,9,'test','test','2018-06-23','2018-06-23','2018-06-23','2018-06-23')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_bill WHERE mid=1")
			rows, err := d.InsertUpBillBatch(c, values)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpSignedAvs(t *testing.T) {
	convey.Convey("ListUpSignedAvs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, last, err := d.ListUpSignedAvs(c, id, limit)
			ctx.Convey("Then err should be nil.ups,last should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(last, convey.ShouldNotBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}
