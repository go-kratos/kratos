package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Matsuri(t *testing.T) {
	Convey("matsuri", t, WithService(func(s *Service) {
		data := s.Matsuri(context.Background(), time.Now())
		So(data, ShouldNotBeNil)
	}))
}

func TestService_View(t *testing.T) {
	Convey("view", t, WithService(func(s *Service) {
		aid := int64(10097666)
		data, err := s.View(context.Background(), aid)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_PageList(t *testing.T) {
	Convey("pagelist", t, WithService(func(s *Service) {
		aid := int64(10097666)
		data, err := s.PageList(context.Background(), aid)
		So(err, ShouldBeNil)
		So(len(data), ShouldBeGreaterThan, 0)
	}))
}

func TestService_VideoShot(t *testing.T) {
	Convey("video shot", t, WithService(func(s *Service) {
		aid := int64(10097666)
		cid := int64(10108404)
		index := true
		data, err := s.VideoShot(context.Background(), aid, cid, index)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}

func TestService_PlayURLToken(t *testing.T) {
	Convey("playurl token", t, WithService(func(s *Service) {
		mid := int64(88895029)
		aid := int64(10097666)
		cid := int64(10108404)
		data, err := s.PlayURLToken(context.Background(), mid, aid, cid)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}
