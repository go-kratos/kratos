package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceMcnRecommendCron(t *testing.T) {
	convey.Convey("McnRecommendCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.McnRecommendCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceDealFailRecommendCron(t *testing.T) {
	convey.Convey("DealFailRecommendCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.DealFailRecommendCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
