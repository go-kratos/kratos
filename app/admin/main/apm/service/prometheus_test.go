package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/apm/model/monitor"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_RPCMonitor(t *testing.T) {
	var (
		c   = context.Background()
		mts []*monitor.Monitor
		err error
	)
	convey.Convey("Prometheus", t, func(ctx convey.C) {
		mts, err = svr.RPCMonitor(c)
		convey.Convey("Then pts should not be nil, err should be nil", func(ctx convey.C) {
			ctx.So(mts, convey.ShouldNotBeEmpty)
			ctx.So(err, convey.ShouldBeNil)
			for _, mt := range mts {
				t.Logf("mt.Interface=%s, mt.Cost=%d, mt.Count=%d", mt.Interface, mt.Cost, mt.Count)
			}
		})
	})

}

func TestService_HttpMonitor(t *testing.T) {
	var (
		c   = context.Background()
		mts []*monitor.Monitor
		err error
	)
	convey.Convey("Prometheus", t, func(ctx convey.C) {
		mts, err = svr.HTTPMonitor(c)
		convey.Convey("Then pts should not be nil, err should be nil", func(ctx convey.C) {
			ctx.So(mts, convey.ShouldNotBeEmpty)
			ctx.So(err, convey.ShouldBeNil)
			for _, mt := range mts {
				t.Logf("mt.Interface=%s, mt.Cost=%d, mt.Count=%d", mt.Interface, mt.Cost, mt.Count)
			}
		})
	})
}

func TestService_TenCent(t *testing.T) {
	var (
		c   = context.Background()
		mts = make([]*monitor.Monitor, 0)
		err error
	)
	convey.Convey("TenCent", t, func(ctx convey.C) {
		mts, err = svr.TenCent(c)
		convey.Convey("Then mts should not be nil, err should be nil", func(ctx convey.C) {
			ctx.So(mts, convey.ShouldNotBeEmpty)
			ctx.So(err, convey.ShouldBeNil)
			for _, mt := range mts {
				t.Logf("mt.AppID=%s, mt.Interface=%s, mt.Count=%d", mt.AppID, mt.Interface, mt.Count)
			}
		})
	})
}

func TestService_KingSoft(t *testing.T) {
	var (
		c   = context.Background()
		mts = make([]*monitor.Monitor, 0)
		err error
	)
	convey.Convey("KingSoft", t, func(ctx convey.C) {
		mts, err = svr.KingSoft(c)
		convey.Convey("Then mts should not be nil, err should be nil", func(ctx convey.C) {
			ctx.So(mts, convey.ShouldNotBeEmpty)
			ctx.So(err, convey.ShouldBeNil)
			for _, mt := range mts {
				t.Logf("mt.AppID=%s, mt.Interface=%s, mt.Count=%d", mt.AppID, mt.Interface, mt.Count)
			}
		})
	})
}

func TestService_DataBus(t *testing.T) {
	var (
		c   = context.Background()
		mts = make([]*monitor.Monitor, 0)
		err error
	)
	convey.Convey("DataBus", t, func(ctx convey.C) {
		mts, err = svr.DataBus(c)
		convey.Convey("Then mts should not be nil, err should be nil", func(ctx convey.C) {
			ctx.So(mts, convey.ShouldNotBeEmpty)
			ctx.So(err, convey.ShouldBeNil)
			for _, mt := range mts {
				t.Logf("mt.AppID=%s, mt.Interface=%s, mt.Count=%d", mt.AppID, mt.Interface, mt.Count)
			}
		})
	})
}
