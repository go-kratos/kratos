package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_OnlineList(t *testing.T) {
	Convey("online list", t, WithService(func(s *Service) {
		data, err := s.OnlineList(context.Background())
		So(err, ShouldBeNil)
		Printf("%v", data)
	}))
}

func TestService_OnlineArchiveCount(t *testing.T) {
	Convey("test online OnlineArchiveCount", t, WithService(func(s *Service) {
		res := s.OnlineArchiveCount(context.Background())
		So(res, ShouldNotBeNil)
	}))
}
