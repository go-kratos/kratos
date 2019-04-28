package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Kv(t *testing.T) {
	Convey("test baidu Kv", t, WithService(func(s *Service) {
		res, err := s.Kv(context.Background())
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_CmtBox(t *testing.T) {
	Convey("test cmtbox", t, WithService(func(s *Service) {
		res, err := s.CmtBox(context.Background(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
