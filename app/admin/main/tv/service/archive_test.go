package service

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_CheckArc(t *testing.T) {
	Convey("TestService_CheckArc Test", t, WithService(func(s *Service) {
		res, err := s.CheckArc(10099763)
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	}))
}

func TestService_ExistArc(t *testing.T) {
	Convey("TestService_ExistArc Test", t, WithService(func(s *Service) {
		res, err := s.ExistArc(12009430)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestService_AddArcs(t *testing.T) {
	Convey("TestService_ArchiveAdd Test", t, WithService(func(s *Service) {
		res, err := s.AddArcs([]int64{
			10099763, 10099764,
		})
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}
