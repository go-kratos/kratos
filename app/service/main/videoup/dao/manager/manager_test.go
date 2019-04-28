package manager

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerUppers(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Uppers", t, func(ctx convey.C) {
		_, err := d.Uppers(c)
		ctx.Convey("Then err should be nil.um should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
