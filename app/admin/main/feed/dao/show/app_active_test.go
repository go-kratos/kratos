package show

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestShow_app_activeFindByID(t *testing.T) {
	convey.Convey("AAFindByID", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.AAFindByID(c, id)
			ctx.Convey("Then err should be nil.active should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
