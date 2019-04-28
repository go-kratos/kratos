package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveDescFormats(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("DescFormats", t, func(ctx convey.C) {
		_, err := d.DescFormats(c)
		ctx.Convey("Then err should be nil.dfs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
