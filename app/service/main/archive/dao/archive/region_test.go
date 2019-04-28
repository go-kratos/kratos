package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveRegionArcs(t *testing.T) {
	var (
		c      = context.TODO()
		rid    = int16(97)
		start  = int(0)
		length = int(10)
	)
	convey.Convey("RegionArcs", t, func(ctx convey.C) {
		_, err := d.RegionArcs(c, rid, start, length)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
