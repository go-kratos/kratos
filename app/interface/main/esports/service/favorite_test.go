package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddFav(t *testing.T) {
	Convey("test service addfav", t, WithService(func(s *Service) {
		err := s.AddFav(context.Background(), 180, 3)
		So(err, ShouldBeNil)
	}))
}

func TestService_DelFav(t *testing.T) {
	Convey("test service delfav", t, WithService(func(s *Service) {
		err := s.DelFav(context.Background(), 180, 0)
		So(err, ShouldBeNil)
	}))
}

func TestService_ListFav(t *testing.T) {
	Convey("test service listfav", t, WithService(func(s *Service) {
		res, count, err := s.ListFav(context.Background(), 10, 10, 1, 10)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan)
		println(count)
	}))
}
