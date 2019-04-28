package history

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_List(t *testing.T) {
	Convey("list", t, WithService(func(s *Service) {
		_, err := s.List(context.TODO(), 27515256, 111, 1, 20, "0", 0)
		So(err, ShouldBeNil)
	}))
}

func TestService_Live(t *testing.T) {
	Convey("live", t, WithService(func(s *Service) {
		_, err := s.Live(context.TODO(), []int64{27515256})
		So(err, ShouldBeNil)
	}))
}

func TestService_LiveList(t *testing.T) {
	Convey("live", t, WithService(func(s *Service) {
		_, err := s.LiveList(context.TODO(), 27515256, 111, 1, 20, "0", 0)
		So(err, ShouldBeNil)
	}))
}
