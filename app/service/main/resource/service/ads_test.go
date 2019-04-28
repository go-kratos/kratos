package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LoadVideoAds(t *testing.T) {
	Convey("load video ads", t, WithService(func(s *Service) {
		err := s.loadVideoAds()
		So(err, ShouldBeNil)
	}))
}

func Test_PasterAPP(t *testing.T) {
	Convey("get nologin paster", t, WithService(func(s *Service) {
		_, err := s.PasterAPP(context.Background(), 0, 1, "15038932", "11", "22222")
		So(err, ShouldBeNil)
	}))
}

func Test_PasterPGC(t *testing.T) {
	Convey("get bangumi paster", t, WithService(func(s *Service) {
		_, err := s.PasterPGC(context.Background(), 2, 0, "43883")
		So(err, ShouldBeNil)
	}))
}

func Test_PasterCID(t *testing.T) {
	Convey("get paster cids", t, WithService(func(s *Service) {
		_, err := s.PasterCID(context.Background())
		So(err, ShouldBeNil)
	}))
}
