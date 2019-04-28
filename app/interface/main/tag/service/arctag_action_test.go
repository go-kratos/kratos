package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Action(t *testing.T) {
	var (
		mid int64 = 14771787
		aid int64 = 4052445
		tid int64 = 10176
		now       = time.Now()
	)
	Convey("testLike service", t, WithService(func(s *Service) {
		testSvc.Like(context.Background(), mid, aid, tid, now)
	}))
	Convey("testHate service", t, WithService(func(s *Service) {
		testSvc.Hate(context.Background(), mid, aid, tid, now)
	}))
}
