package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UserLikes(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		res, err := s.UserLikes(c, "danmu", 1, 1, 2)
		So(err, ShouldBeNil)
		So(len(res), ShouldEqual, 1)
	}))
}

func Test_UserTotalLike(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		_, err := s.UserTotalLike(c, "danmu", 1, 1, 2)
		So(err, ShouldBeNil)
	}))
}

func Test_UserDislikes(t *testing.T) {
	Convey("get data", t, WithService(func(s *Service) {
		_, err := s.UserDislikes(c, "danmu", 1, 1, 2)
		So(err, ShouldBeNil)
	}))
}
