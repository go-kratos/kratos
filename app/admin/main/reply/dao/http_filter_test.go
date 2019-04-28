package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoFilterContents(t *testing.T) {
	convey.Convey("FilterContents", t, func(ctx convey.C) {
		var (
			rpMaps map[int64]string
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.FilterContents(context.Background(), rpMaps)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
