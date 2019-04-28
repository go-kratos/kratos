package service

import (
	"context"
	"testing"

	"go-common/app/service/main/videoup/model/archive"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddByUGC(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("AddByUGC", t, WithService(func(s *Service) {
		_, err := svr.AddByUGC(c, ap)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_EditByUGC(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("EditByUGC", t, WithService(func(s *Service) {
		err := svr.EditByUGC(c, ap)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_EditMission(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcMissionParam{AID: 5463716, MID: 15555180, MissionID: 1211, Tag: "11,22,33"}
	)
	Convey("EditByUGC", t, WithService(func(s *Service) {
		err := svr.EditMissionByUGC(c, ap)
		spew.Dump(err)
		So(err, ShouldNotBeNil)
	}))
}
