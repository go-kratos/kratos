package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeWebDataCnt(t *testing.T) {
	convey.Convey("WebDataCnt", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			vid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.WebDataCnt(c, vid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeWebDataList(t *testing.T) {
	convey.Convey("WebDataList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			vid    = int64(36)
			offset = int(1)
			limit  = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			list, err := d.WebDataList(c, vid, offset, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}
