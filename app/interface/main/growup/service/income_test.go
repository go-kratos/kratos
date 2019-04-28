package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArchiveIncome(t *testing.T) {
	var (
		mid        = int64(1011)
		typ        = 0
		page, size = 0, 10
		all        = 0
	)
	Convey("ArchiveIncome", t, WithService(func(s *Service) {
		_, err := s.ArchiveIncome(context.Background(), mid, typ, page, size, all)
		So(err, ShouldBeNil)
	}))
}

func Test_UpSummary(t *testing.T) {
	Convey("UpSummary", t, WithService(func(s *Service) {
		var mid int64 = 1011
		_, err := s.UpSummary(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}
