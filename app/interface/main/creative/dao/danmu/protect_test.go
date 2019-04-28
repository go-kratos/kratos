package danmu

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDanmuProtectApplyList(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(2089809)
		page   = int64(0)
		aidStr = "1"
		sort   = "ctime"
		ip     = "127.0.0.1"
	)
	convey.Convey("ProtectApplyList", t, func(ctx convey.C) {
		result, err := d.ProtectApplyList(c, mid, page, aidStr, sort, ip)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuProtectApplyVideoList(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		ip  = "127.0.0.1"
	)
	convey.Convey("ProtectApplyVideoList", t, func(ctx convey.C) {
		result, err := d.ProtectApplyVideoList(c, mid, ip)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuProtectOper(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(2089809)
		status = int64(0)
		ids    = "1,2"
		ip     = "127.0.0.1"
	)
	convey.Convey("ProtectOper", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.dmProtectApplyStatusURL).Reply(200).JSON(`{"code":20043,"data":""}`)
		err := d.ProtectOper(c, mid, status, ids, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
