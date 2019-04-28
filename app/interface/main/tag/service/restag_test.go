package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ResTag(t *testing.T) {
	var (
		name           = "妖艳"
		now            = time.Now()
		tagNames       = []string{"AV", "妖艳"}
		oids           = []int64{78, 75, 77}
		oid      int64 = 78
		mid      int64 = 14771787
		typ      int8  = 1
	)
	Convey("testResTags service", t, WithService(func(s *Service) {
		testSvc.ResTags(context.Background(), oids, mid, typ)
	}))
	Convey("testUpResBind service", t, WithService(func(s *Service) {
		testSvc.UpResBind(context.Background(), oid, mid, tagNames, typ, now)
	}))
	Convey("testResAdminBind service", t, WithService(func(s *Service) {
		testSvc.ResAdminBind(context.Background(), oid, mid, tagNames, typ, now)
	}))
	Convey("testCheckName service", t, WithService(func(s *Service) {
		testSvc.CheckName(name)
	}))
}
