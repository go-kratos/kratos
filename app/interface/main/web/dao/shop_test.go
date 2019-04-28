package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoShopInfo(t *testing.T) {
	convey.Convey("ShopInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515399)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.ShopInfo(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}
