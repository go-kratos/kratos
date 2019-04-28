package manager

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestManagerArcReason(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(2333)
	)
	convey.Convey("ArcReason", t, func(ctx convey.C) {
		_, err := d.ArcReason(c, aid)
		ctx.Convey("Then err should be nil.tagID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
