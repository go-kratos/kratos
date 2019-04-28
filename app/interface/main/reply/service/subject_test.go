package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubject(t *testing.T) {
	c := context.Background()
	Convey("subject test", t, WithService(func(s *Service) {
		_, err := s.Subject(c, 1, 1)
		So(err, ShouldBeNil)
	}))
	Convey("subject test", t, WithService(func(s *Service) {
		_, err := s.AdminGetSubject(c, 1, 1)
		So(err, ShouldBeNil)
	}))
	Convey("subject test", t, WithService(func(s *Service) {
		err := s.AdminSubjectMid(c, 1, 1, 1, 1, "test")
		So(err, ShouldBeNil)
	}))
	Convey("subject test", t, WithService(func(s *Service) {
		err := s.AdminSubRegist(c, 1, 1, 1, 1, "test")
		So(err, ShouldBeNil)
	}))
	Convey("subject test", t, WithService(func(s *Service) {
		_, err := s.GetSubject(c, 1, 1)
		So(err, ShouldBeNil)
	}))
}
