package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpMcnDataSummaryCron(t *testing.T) {
	convey.Convey("UpMcnDataSummaryCron", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.UpMcnDataSummaryCron()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
