package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_ViewPoints(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("ViewPoint", t, func(ctx convey.C) {
		vp, err := d.ViewPoint(c, 10106351, 10126396)
		ctx.Convey("", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vp, convey.ShouldNotBeNil)
		})
	})
}
