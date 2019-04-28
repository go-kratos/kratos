package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPingMc(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
