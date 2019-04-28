package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveTypes(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Types", t, func(ctx convey.C) {
		types, err := d.Types(c)
		ctx.Convey("Then err should be nil.types should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(types, convey.ShouldNotBeNil)
		})
	})
}
