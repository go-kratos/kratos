package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/mcn/model/mcnmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceCmdReloadRankCache(t *testing.T) {
	convey.Convey("CmdReloadRankCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &mcnmodel.CmdReloadRank{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.CmdReloadRankCache(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
