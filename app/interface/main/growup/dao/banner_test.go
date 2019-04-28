package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBanner(t *testing.T) {
	convey.Convey("Banner", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			no = int64(1541409804)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			b, err := d.Banner(c, no)
			Exec(c, "INSERT INTO banner(id, start_at, end_at) VALUES(1000, '2018-01-01', '2019-01-01')")
			ctx.Convey("Then err should be nil.b should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(b, convey.ShouldNotBeNil)
			})
		})
	})
}
