package favorite

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIsFavDefault(t *testing.T) {
	convey.Convey("IsFavDefault", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IsFavDefault(c, 1, 1)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestIsFav(t *testing.T) {
	convey.Convey("IsFav", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.IsFav(c, 1, 1)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestAddFav(t *testing.T) {
	convey.Convey("AddFav", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.AddFav(c, 1, 1)
		})
	})
}
