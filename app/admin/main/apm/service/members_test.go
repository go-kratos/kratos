package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/monitor"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_Members(t *testing.T) {
	var (
		c   = context.Background()
		mt  *monitor.Monitor
		err error
	)
	convey.Convey("Members", t, func(ctx convey.C) {
		mt, err = svr.Members(c)
		ctx.So(mt, convey.ShouldNotBeEmpty)
		ctx.So(err, convey.ShouldBeNil)
		t.Logf("AppID=%s, Interface=%s, Count=%d", mt.AppID, mt.Interface, mt.Count)
	})
}
