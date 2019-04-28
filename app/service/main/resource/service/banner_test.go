package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LoadBannerCahce(t *testing.T) {
	Convey("load banner cache", t, WithService(func(s *Service) {
		err := s.loadBannerCahce()
		So(err, ShouldBeNil)
	}))
}

func Test_Banners(t *testing.T) {
	Convey("get app banner", t, WithService(func(s *Service) {
		res := s.Banners(context.Background(), 1, 6190, 0, 182504479, "457", "", "218.4.147.222", "222", "wifi", "iphone", "phone", "", "", "", true)
		So(res, ShouldNotBeNil)
	}))
}

func Test_CPMBanners(t *testing.T) {
	Convey("get cpm by api", t, WithService(func(s *Service) {
		res := s.cpmBanners(context.Background(), 0, 182504479, 6190, "457", "iphone", "phone", "222", "wifi", "218.4.147.222", "", "")
		So(res, ShouldNotBeEmpty)
	}))
}
