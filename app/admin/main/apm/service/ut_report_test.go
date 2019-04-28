package service

import (
	"context"
	"go-common/app/admin/main/apm/dao"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestServiceWechatReport(t *testing.T) {
	convey.Convey("WechatReport", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mrid = int64(10760)
			cid  = "8d2f1b49661c7089e2b595eafff326033a138c23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(svr.dao), "SendWechatToGroup", func(_ *dao.Dao, _ context.Context, _ string, _ string) error {
				return nil
			})
			defer guard.Unpatch()
			err := svr.WechatReport(c, mrid, cid, "xxx", "master")
			ctx.Convey("Than err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceRankWechatReport(t *testing.T) {
	convey.Convey("RankWechatReport", t, func(ctx convey.C) {
		c := context.Background()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(svr.dao), "SendWechatToGroup", func(_ *dao.Dao, _ context.Context, _ string, _ string) error {
				return nil
			})
			defer guard.Unpatch()
			err := svr.RankWechatReport(c)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceSummaryWechatReport(t *testing.T) {
	convey.Convey("SummaryWechatReport", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(svr.dao), "SendWechatToGroup", func(_ *dao.Dao, _ context.Context, _ string, _ string) error {
				return nil
			})
			defer guard.Unpatch()
			svr.dao.SetAppCovCache(c)
			err := svr.SummaryWechatReport(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				// t.Logf("\n%s", msg)
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
