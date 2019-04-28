package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ResTag(t *testing.T) {
	Convey("Test_ResTag info  service ", t, WithService(func(s *Service) {
		var (
			ip     string
			c            = context.Background()
			mid    int64 = 14771787
			oid    int64 = 1845325
			typ    int32 = 3
			ps, pn       = 2, 3
		)
		s.ResTags(c, oid, typ, mid)
		s.ResTagLog(c, oid, typ, mid, pn, ps, ip)
	}))
}
