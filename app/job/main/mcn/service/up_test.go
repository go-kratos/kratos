package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpMcnUpStateCron(t *testing.T) {
	convey.Convey("UpMcnUpStateCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.UpMcnUpStateCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
