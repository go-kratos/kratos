package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Business(t *testing.T) {
	c := context.Background()
	Convey("test service business", c, WithService(func(s *Service) {
		Convey("test service add business", c, WithService(func(s *Service) {
			_, err := s.AddBusiness(c, -3, "abc", "abc", "abc", "abc")
			So(err, ShouldBeNil)
		}))
		Convey("test service list business", c, WithService(func(s *Service) {
			_, err := s.ListBusiness(c, 0)
			So(err, ShouldBeNil)
		}))
		Convey("test service update business", c, WithService(func(s *Service) {
			_, err := s.UpBusiness(c, "test", "test", "test", "abc", -3)
			So(err, ShouldBeNil)
		}))
		Convey("test service update business state", c, WithService(func(s *Service) {
			_, err := s.UpBusinessState(c, 1, -3)
			So(err, ShouldBeNil)
		}))
	}))
}
