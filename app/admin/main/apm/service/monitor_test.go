package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/monitor"
	"go-common/library/log"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_AddMonitor(t *testing.T) {
	var (
		c   = context.Background()
		err error
	)
	convey.Convey("Test AddMonitor", t, func(ctx convey.C) {
		err = svr.AddMonitor(c)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestService_PrometheusList(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Test PrometheusList", t, func(ctx convey.C) {
		ret, err := svr.PrometheusList(c, "main.account.member-service", "xx", "count")
		ctx.So(ret, convey.ShouldNotBeEmpty)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestService_BroadCastList(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("Test BroadCastList", t, func(ctx convey.C) {
		ret, err := svr.BroadCastList(c)
		ctx.So(ret, convey.ShouldNotBeEmpty)
		ctx.So(err, convey.ShouldBeNil)
	})
}

func TestService_Times(t *testing.T) {
	mt := &monitor.Monitor{}
	if err := svr.DB.Where("app_id = ?", "main.app-svr.app-feed-http").First(mt).Error; err != nil {
		log.Error("s.Prometheus query first error(%v)", err)
		return
	}
	convey.Convey("Test Times", t, func(ctx convey.C) {
		ret := svr.times(mt.MTime)
		t.Logf("times=%v", ret)
	})
}
