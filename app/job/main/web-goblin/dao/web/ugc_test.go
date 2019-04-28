package web

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestWebUgcSearch(t *testing.T) {
	convey.Convey("UgcSearch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data map[string]interface{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UgcSearch(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
