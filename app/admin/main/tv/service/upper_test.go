package service

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddMids(t *testing.T) {
	Convey("TestService_CheckMids Test", t, WithService(func(s *Service) {
		res, err := s.AddMids([]int64{
			1, 2, 3, 777777777777, 999,
		})
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestService_ImportMids(t *testing.T) {
	Convey("TestService_ImportMids Test", t, WithService(func(s *Service) {
		res, err := s.ImportMids([]int64{
			1, 2, 3, 777777777777, 999,
		})
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestService_DelMid(t *testing.T) {
	Convey("TestService_DelMid Test", t, WithService(func(s *Service) {
		err := s.DelMid(2)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpList(t *testing.T) {
	Convey("TestService_UpList Test", t, WithService(func(s *Service) {
		res, err := s.UpList(2, 1, "", 0)
		So(err, ShouldBeNil)
		So(len(res.Items), ShouldBeGreaterThan, 0)
	}))
}
