package report

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReportReport(t *testing.T) {
	convey.Convey("Report", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Report(c, "archive_click")
			println(err)
		})
	})
}

func TestReportCheckJob(t *testing.T) {
	convey.Convey("CheckJob", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			urls = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.CheckJob(c, urls)
		})
	})
}
