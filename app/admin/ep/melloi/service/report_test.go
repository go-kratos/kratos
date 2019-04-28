package service

import (
	"testing"

	"go-common/app/admin/ep/melloi/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Report(t *testing.T) {

	reportSummaryC := model.ReportSummary{
		TestName: "testGrpc",
	}
	Convey("query reportGraph", t, func() {
		var strs = []string{"hoopchina1534255529"}
		regraphs, _ := s.QueryReGraph(strs)
		So(len(regraphs), ShouldBeGreaterThan, 0)
	})

	Convey("count query report summarys", t, func() {
		count, _ := s.CountQueryReportSummarys(&reportSummaryC)
		So(count, ShouldBeGreaterThan, 0)
	})

	Convey("add reportGraph", t, func() {
		//reportSummarys:= s.addReportGraph(&reportSummaryC,1,10)
		//So(len(reportSummarys), ShouldBeGreaterThan, 0)
	})
}
