package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Report(t *testing.T) {
	Convey("Test_Report    service ", t, WithService(func(s *Service) {
		var (
			ip  string
			c         = context.Background()
			mid int64 = 14771787
			oid int64 = 1845325
			typ int32 = 3
		)
		s.ReportAction(c, oid, 11, mid, typ, 1, 2, 0, "", ip)
	}))
}
