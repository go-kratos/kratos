package gorm

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"time"
)

func TestGormReportTaskFlow(t *testing.T) {
	convey.Convey("ReportTaskFlow", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			metas, _, err := d.ReportTaskMetas(c, "", "", 1, 1, []int64{}, map[int64]string{}, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				for _, meta := range metas {
					fmt.Println("meta:", meta)
				}
			})
		})
	})
}

func TestDao_TaskReports(t *testing.T) {
	time.Time{}.Unix()
	t.Logf("%s", time.Duration(time.Duration(3661)*time.Second).String())
	return
	convey.Convey("TaskReports", t, func(ctx convey.C) {
		res, err := d.TaskReports(cntx, 1, 1, []int8{2, 3, 4}, "2019-01-14", "2019-01-15")
		ctx.So(err, convey.ShouldBeNil)
		for _, item := range res {
			t.Logf("item(%+v)", item)
		}
	})
}
