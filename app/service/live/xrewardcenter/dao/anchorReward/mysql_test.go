package anchorReward

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAnchorRewardRawRewardConf(t *testing.T) {
	convey.Convey("RawRewardConf", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.RawRewardConf(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
