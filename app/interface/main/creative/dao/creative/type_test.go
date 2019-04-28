package creative

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCreativeTypes(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Types", t, func(ctx convey.C) {
		tops, langs, typeMap, err := d.Types(c)
		ctx.Convey("Then err should be nil.tops,langs,typeMap should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(typeMap, convey.ShouldNotBeNil)
			ctx.So(langs, convey.ShouldNotBeNil)
			ctx.So(tops, convey.ShouldNotBeNil)
		})
	})
}
