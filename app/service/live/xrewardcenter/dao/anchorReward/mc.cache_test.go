package anchorReward

import (
	"context"
	model "go-common/app/service/live/xrewardcenter/model/anchorTask"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAnchorRewardCacheRewardConf(t *testing.T) {
	convey.Convey("CacheRewardConf", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheRewardConf(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestAnchorRewardAddCacheRewardConf(t *testing.T) {
	convey.Convey("AddCacheRewardConf", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(0)
			val = &model.AnchorRewardConf{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheRewardConf(c, id, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
