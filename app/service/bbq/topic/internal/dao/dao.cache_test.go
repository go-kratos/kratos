package dao

import (
	"context"
	"go-common/library/log"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoVideoExtension(t *testing.T) {
	convey.Convey("VideoExtension", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.VideoExtension(c, keys)
			log.V(1).Infow(c, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTopicInfo(t *testing.T) {
	convey.Convey("TopicInfo", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.TopicInfo(c, keys)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
