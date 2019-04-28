package service

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model"
)

func TestService_ReportTaskSubmit(t *testing.T) {
	convey.Convey("ReportTaskSubmit", t, func(ctx convey.C) {
		pm := &model.OptReportSubmit{
			BizID: 2,
		}

		res, err := s.ReportTaskSubmit(cntx, pm)
		ctx.So(err, convey.ShouldBeNil)
		t.Logf("res.header(%+v)", res.Header)
		for _, list := range res.Rows {
			for _, item := range list {
				t.Logf("item(%+v)", item)
			}
		}
	})
}
