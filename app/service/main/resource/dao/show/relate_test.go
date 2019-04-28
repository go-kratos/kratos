package show

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRelate(t *testing.T) {
	convey.Convey("Relate", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			_, err := d.Relate(context.Background(), time.Now())
			ctx.Convey("Error should be nil, res should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
