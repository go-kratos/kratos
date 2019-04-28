package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestStaff(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(4052032)
	)
	convey.Convey("Staff", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			_, err := d.Staff(c, aid)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
