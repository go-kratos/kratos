package service

import (
	"context"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAction(t *testing.T) {
	c := context.Background()
	Convey("action set", t, WithService(func(s *Service) {
		err := s.UpActionLike(c, 894717392, 5464686, 23, 1, 32, "test")
		So(err, ShouldBeNil)
		like, hate, err := s.ActionCount(c, 894717392, 5464686, 23, 1)
		So(err, ShouldBeNil)
		So(hate, ShouldEqual, 0)
		So(like, ShouldEqual, 32)
	}))
}
