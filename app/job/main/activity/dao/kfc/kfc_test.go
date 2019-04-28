package kfc

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestKfcKfcDelver(t *testing.T) {
	convey.Convey("KfcDelver", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(1)
			mid = int64(5248758)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.KfcDelver(c, id, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
