package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoQueryTaskLog(t *testing.T) {
	convey.Convey("QueryTaskLog", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			queryStmt = ""
			sort      = "ctime"
			offset    = int(1)
			limit     = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.QueryTaskLog(c, queryStmt, sort, offset, limit)
			ctx.Convey("Then err should be nil.taskLogs,count should not be nil.", func(ctx convey.C) {

			})
		})
	})
}
