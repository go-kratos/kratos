package manager

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerUpSpecial(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("UpSpecial", t, func(ctx convey.C) {
		_, err := d.UpSpecial(c)
		ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
