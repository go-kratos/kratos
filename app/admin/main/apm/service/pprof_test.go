package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/pprof"

	"github.com/smartystreets/goconvey/convey"
)

var (
	text = `
{
    "title": "主站 HTTP_SERVER 错误率过高(主告警条件)",
    "tags": {
        "app": "account.service.member",
        "code": "-404",
        "exported_job": "caster_app_metrics",
        "method": "x/v2/view"
    }
}`
)

func TestService_ActiveWarning(t *testing.T) {
	convey.Convey("ActiveWarning", t, func() {
		err := svr.ActiveWarning(context.Background(), text)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestService_Pprof(t *testing.T) {
	var (
		err error
		req = &pprof.Params{
			AppID:   "account.service.member",
			Kind:    1,
			SvgName: "4zf56-1539587841",
		}
		pws = make([]*pprof.Warn, 0)
	)
	convey.Convey("PprofWarn", t, func() {
		pws, err = svr.PprofWarn(context.Background(), req)
		convey.So(err, convey.ShouldBeNil)
		convey.So(pws, convey.ShouldNotBeEmpty)
		for _, pw := range pws {
			t.Logf("pw.Kind=%d, pw.URL=%s", pw.Kind, pw.URL)
		}
	})
}
