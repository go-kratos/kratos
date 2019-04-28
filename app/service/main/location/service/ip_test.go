package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Info(t *testing.T) {
	Convey("ip info", t, WithService(func(s *Service) {
		res, err := s.Info(context.Background(), "211.139.80.6")
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func Test_Infos(t *testing.T) {
	Convey("ips info", t, WithService(func(s *Service) {
		res, err := s.Infos(context.Background(), []string{"211.139.80.6", "64.233.173.24"})
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_InfoComplete(t *testing.T) {
	Convey("ip complete info", t, WithService(func(s *Service) {
		res, err := s.InfoComplete(context.Background(), "211.139.80.6")
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func Test_InfosComplete(t *testing.T) {
	Convey("ips complete info", t, WithService(func(s *Service) {
		res, err := s.InfosComplete(context.Background(), []string{"211.139.80.6", "64.233.173.24"})
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
