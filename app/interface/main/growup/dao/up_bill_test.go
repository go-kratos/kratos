package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpBill(t *testing.T) {
	convey.Convey("GetUpBill", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1001)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_bill(mid) VALUES(1001)")
			up, err := d.GetUpBill(c, mid)
			ctx.Convey("Then err should be nil.up should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(up, convey.ShouldNotBeNil)
			})
		})
	})
}
