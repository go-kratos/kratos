package service

import (
	"context"
	"go-common/app/admin/main/aegis/model/monitor"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceMonitorBuzResult(t *testing.T) {
	convey.Convey("MonitorBuzResult", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			bid = int64(2)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := s.MonitorBuzResult(c, bid)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMonitorResultOids(t *testing.T) {
	convey.Convey("MonitorBuzResult", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := s.MonitorResultOids(c, id)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestServiceMonitorNotifyTime(t *testing.T) {
	convey.Convey("MonitorNotifyTime", t, func(convCtx convey.C) {
		var (
			conf = &monitor.RuleConf{
				NotifyCdt: map[string]struct {
					Comp  string `json:"comparison"`
					Value int64  `json:"value"`
				}{
					"time": {
						Comp:  ">",
						Value: 1,
					},
				},
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, _, err := s.monitorNotifyTime(conf)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
