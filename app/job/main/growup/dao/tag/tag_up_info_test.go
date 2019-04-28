package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagGetTagUpInfoByTag(t *testing.T) {
	convey.Convey("GetTagUpInfoByTag", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tags   = []int64{1, 2, 3}
			from   = int(0)
			limit  = int(10)
			tagMID map[int64][]int64
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			tagMID = make(map[int64][]int64)
			count, err := d.GetTagUpInfoByTag(c, tags, from, limit, tagMID)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
