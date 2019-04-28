package service

import (
	"testing"
	"time"

	"context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ShortCache(t *testing.T) {
	Convey("ShortCache", t, WithService(func(s *Service) {
		_, err := s.ShortCache(context.TODO(), "http://b23.tv/EbUzmu")
		So(err, ShouldBeNil)
	}))
}

func TestService_ShortByID(t *testing.T) {
	Convey("ShortByID", t, WithService(func(s *Service) {
		_, err := s.ShortByID(context.TODO(), 1)
		So(err, ShouldBeNil)
	}))
}

func TestService_Add(t *testing.T) {
	Convey("Add", t, WithService(func(s *Service) {
		_, err := s.Add(context.TODO(), 279, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}

func TestService_ShortUpdate(t *testing.T) {
	Convey("ShortUpdate", t, WithService(func(s *Service) {
		err := s.ShortUpdate(context.TODO(), 1, 279, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}

func TestService_ShortDel(t *testing.T) {
	Convey("ShortDel", t, WithService(func(s *Service) {
		err := s.ShortDel(context.TODO(), 1, 279, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestService_ShortCount(t *testing.T) {
	Convey("ShortCount", t, WithService(func(s *Service) {
		_, err := s.ShortCount(context.TODO(), 279, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}

func TestService_ShortLimit(t *testing.T) {
	Convey("ShortLimit", t, WithService(func(s *Service) {
		_, err := s.ShortLimit(context.TODO(), 1, 20, 279, "http://www.baidu.com")
		So(err, ShouldBeNil)
	}))
}
