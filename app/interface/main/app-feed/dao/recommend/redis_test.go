package recommend

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPositionCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.PositionCache(c, mid)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestAddPositionCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		pos = int(1)
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.AddPositionCache(c, mid, pos)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
