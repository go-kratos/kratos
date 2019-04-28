package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArchiveStatis(t *testing.T) {
	Convey("ArchiveStatis", t, WithService(func(s *Service) {
		categoryID := []int64{}
		groupType := 1
		var fromTime, endTime int64 = 1524240000000, 1526832000000
		_, err := s.ArchiveStatis(context.Background(), categoryID, 0, groupType, fromTime, endTime)
		So(err, ShouldBeNil)
	}))
}

func Test_ArchiveSection(t *testing.T) {
	Convey("ArchiveSection", t, WithService(func(s *Service) {
		categoryID := []int64{}
		groupType := 1
		var fromTime, endTime int64 = 1524240000000, 1526832000000
		_, err := s.ArchiveSection(context.Background(), categoryID, 0, groupType, fromTime, endTime)
		So(err, ShouldBeNil)
	}))
}

func Test_ArchiveDetail(t *testing.T) {
	Convey("ArchiveDetail", t, WithService(func(s *Service) {
		mid := int64(1)
		groupType := 1
		var fromTime, endTime int64 = 1524240000000, 1526832000000
		_, err := s.ArchiveDetail(context.Background(), mid, 0, groupType, fromTime, endTime)
		So(err, ShouldBeNil)
	}))
}

func Test_ArchiveTop(t *testing.T) {
	Convey("ArchiveTop", t, WithService(func(s *Service) {
		avIDs := []int64{}
		groupType := 1
		var fromTime, endTime int64 = 1524240000000, 1526832000000
		from, limit := 0, 10
		_, _, err := s.ArchiveTop(context.Background(), avIDs, 0, groupType, fromTime, endTime, from, limit)
		So(err, ShouldBeNil)
	}))
}
