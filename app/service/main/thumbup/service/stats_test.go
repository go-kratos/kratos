package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Stats(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		_, err := s.Stats(c, "danmu", 1, []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func Test_StatsWithLike(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		_, err := s.StatsWithLike(c, "danmu", 1, 1, []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func Test_UpdateCount(t *testing.T) {
	Convey("work", t, WithService(func(s *Service) {
		err := s.UpdateCount(c, "danmu", 1, 3, 1, 1, "", "sam")
		So(err, ShouldBeNil)
	}))
}
