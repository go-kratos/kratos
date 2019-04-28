package timemachine

import (
	"context"
	"go-common/app/interface/main/activity/model/timemachine"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTimemachineAddCacheTimemachine(t *testing.T) {
	convey.Convey("AddCacheTimemachine", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &timemachine.Item{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheTimemachine(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestTimemachineCacheTimemachine(t *testing.T) {
	convey.Convey("CacheTimemachine", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheTimemachine(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
