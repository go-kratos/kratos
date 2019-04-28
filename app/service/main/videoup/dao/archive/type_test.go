package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveTypeMapping(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("TypeMapping", t, func(ctx convey.C) {
		_, err := d.TypeMapping(c)
		ctx.Convey("Then err should be nil.tmap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
