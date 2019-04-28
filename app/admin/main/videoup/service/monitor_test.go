package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/model/monitor"
	"testing"
)

func TestService_MonitorResult(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("MonitorResult", t, WithService(func(s *Service) {
		p := &monitor.RuleResultParams{
			Type:     1,
			Business: 1,
		}
		data, err := svr.MonitorResult(c, p)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_MonitorCheckVideoStatus(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("MonitorCheckVideoStatus", t, WithService(func(s *Service) {
		err := svr.MonitorCheckVideoStatus(c, 1, 1)
		So(err, ShouldBeNil)
	}))
}
