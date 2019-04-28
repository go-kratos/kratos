package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveAddit(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
	)
	convey.Convey("Addit", t, func(ctx convey.C) {
		_, err := d.Addit(c, aid)
		ctx.Convey("Then err should be nil.addit should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
