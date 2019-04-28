package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_PorderCfgList(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("PorderCfgList", t, WithService(func(s *Service) {
		_, err := svr.PorderCfgList(c)
		So(err, ShouldBeNil)
	}))
}

func TestService_PorderArcList(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("PorderArcList", t, WithService(func(s *Service) {
		_, err := svr.PorderArcList(c, "", "")
		So(err, ShouldNotBeNil)
	}))
}
