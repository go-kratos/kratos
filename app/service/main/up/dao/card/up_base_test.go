package card

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCardListUpBase(t *testing.T) {
	convey.Convey("ListUpBase", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			size   = int(3)
			lastID = int64(1)
			where  = "AND activity = 1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			idMids, err := d.ListUpBase(c, size, lastID, where)
			ctx.Convey("Then err should be nil.idMids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(idMids, convey.ShouldNotBeNil)
			})
		})
	})
}
