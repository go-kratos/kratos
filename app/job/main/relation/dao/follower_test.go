package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDelFollowerCache(t *testing.T) {
	var (
		fid = int64(0)
	)
	convey.Convey("DelFollowerCache", t, func(ctx convey.C) {
		err := d.DelFollowerCache(fid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
