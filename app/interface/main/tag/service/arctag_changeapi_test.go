package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ArcChangeApi(t *testing.T) {
	var (
		mid    int64 = 14771787
		aid    int64 = 4052445
		tNames       = []string{"AV", "美女"}
		now          = time.Now()
	)
	Convey("testUpArcBind service", t, WithService(func(s *Service) {
		testSvc.UpArcBind(context.Background(), aid, mid, tNames, tNames, now)
	}))
	Convey("testArcAdminBind service", t, WithService(func(s *Service) {
		testSvc.ArcAdminBind(context.Background(), aid, mid, tNames, tNames, now)
	}))
}
