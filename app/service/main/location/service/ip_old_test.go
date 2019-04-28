package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TmpInfo(t *testing.T) {
	Convey("ip tmp info", t, WithService(func(s *Service) {
		res, err := s.TmpInfo("211.139.80.6")
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func Test_TmpInfos(t *testing.T) {
	Convey("ips tmp info", t, WithService(func(s *Service) {
		res, err := s.TmpInfos(context.Background(), []string{"211.139.80.6", "64.233.173.24"}...)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func Test_TmpInfo2(t *testing.T) {
	Convey("ip tmp info", t, WithService(func(s *Service) {
		res, err := s.TmpInfo2(context.Background(), "211.139.80.6")
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}
