package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceHandleBuildContributors(t *testing.T) {
	convey.Convey("HandleBuildContributors", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.HandleBuildContributors(c, repoTest)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
