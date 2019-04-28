package service

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_EpResult(t *testing.T) {
	Convey("EPResult Test", t, WithService(func(s *Service) {
		var req = url.Values{}
		pager, err := s.EpResult(req, 1, 1)
		So(err, ShouldBeNil)
		So(len(pager.Items), ShouldBeGreaterThan, 0)
	}))
}

func TestService_SeasonResult(t *testing.T) {
	Convey("SeasonResult Test", t, WithService(func(s *Service) {
		var req = url.Values{}
		pager, err := s.SeasonResult(req, 1, 1)
		So(err, ShouldBeNil)
		So(len(pager.Items), ShouldBeGreaterThan, 0)
	}))
}
