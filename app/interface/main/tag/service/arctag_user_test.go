package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArcTag(t *testing.T) {
	var (
		ps           = 1
		pn           = 10
		name         = "美丽"
		now          = time.Now()
		reason int8  = 1
		mid    int64 = 14771787
		rptMid int64 = 14771787
		aid    int64 = 4052445
		tid    int64 = 10176
	)
	Convey("testArcTags service", t, WithService(func(s *Service) {
		testSvc.ArcTags(context.Background(), aid, mid)
	}))
	Convey("testLogs service", t, WithService(func(s *Service) {
		testSvc.Logs(context.Background(), aid, mid, pn, ps)
	}))
	Convey("testAdd service", t, WithService(func(s *Service) {
		testSvc.Add(context.Background(), aid, mid, name, now)
	}))
	Convey("testDel service", t, WithService(func(s *Service) {
		testSvc.Del(context.Background(), aid, mid, tid, now)
	}))
	Convey("testAddReport service", t, WithService(func(s *Service) {
		testSvc.AddReport(context.Background(), aid, tid, rptMid, reason, now)
	}))
}
