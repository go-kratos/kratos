package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Action(t *testing.T) {
	Convey("Test_Action service ", t, WithService(func(s *Service) {
		var (
			ip  string
			c         = context.Background()
			tid int64 = 2090
			mid int64 = 14771787
			oid int64 = 1845325
			typ int32 = 3
		)
		s.Like(c, mid, oid, tid, typ, ip)
		s.Hate(c, mid, oid, tid, typ, ip)
	}))
}
