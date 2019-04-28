package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/app-player/model/archive"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddArchiveCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(1)
		arc = &archive.Info{Aid: 1}
	)
	convey.Convey("AddArchiveCache", t, func(ctx convey.C) {
		err := d.AddArchiveCache(c, aid, arc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
