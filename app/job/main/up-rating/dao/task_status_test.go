package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertTaskStatus(t *testing.T) {
	convey.Convey("InsertTaskStatus", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = "(1,2,'2018-06-01','test')"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.InsertTaskStatus(c, val)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
