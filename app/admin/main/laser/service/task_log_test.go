package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceQueryTaskLog(t *testing.T) {
	convey.Convey("QueryTaskLog", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			taskID    = int64(0)
			platform  = int(0)
			taskState = int(0)
			sortBy    = ""
			pageNo    = int(0)
			pageSize  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s.QueryTaskLog(c, mid, taskID, platform, taskState, sortBy, pageNo, pageSize)
			ctx.Convey("Then err should be nil.logs,count should not be nil.", func(ctx convey.C) {

			})
		})
	})
}
