package like

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddExtend(t *testing.T) {
	convey.Convey("ipRequestKey", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = "(2355,10)"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.AddExtend(c, query)
			ctx.Convey("Then err should be nil.likes should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%d", res)
			})
		})
	})
}
