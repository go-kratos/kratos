package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRelateAids(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		_, err := d.RelateAids(c, aid)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
