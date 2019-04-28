package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpMcnSignStateCron(t *testing.T) {
	convey.Convey("UpMcnSignStateCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.UpMcnSignStateCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceUpExpirePayCron(t *testing.T) {
	convey.Convey("UpExpirePayCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.UpExpirePayCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
