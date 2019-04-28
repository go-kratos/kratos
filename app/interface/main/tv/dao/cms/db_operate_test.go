package cms

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsUnshelveArcs(t *testing.T) {
	var (
		c   = context.Background()
		ids = []int64{123}
	)
	convey.Convey("UnshelveArcs", t, func(ctx convey.C) {
		err := d.UnshelveArcs(c, ids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
