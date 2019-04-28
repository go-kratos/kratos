package search

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchLogPaginate(t *testing.T) {
	convey.Convey("LogPaginate", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			oid       = int64(0)
			tp        = int(0)
			states    = []int64{}
			curPage   = int(0)
			pageSize  = int(0)
			startTime = ""
			now       = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			logs, replyCount, reportCount, pageCount, total, err := d.LogPaginate(c, oid, tp, states, curPage, pageSize, startTime, now)
			ctx.Convey("Then err should be nil.logs,replyCount,reportCount,pageCount,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.So(pageCount, convey.ShouldNotBeNil)
				ctx.So(reportCount, convey.ShouldNotBeNil)
				ctx.So(replyCount, convey.ShouldNotBeNil)
				ctx.So(logs, convey.ShouldBeNil)
			})
		})
	})
}
