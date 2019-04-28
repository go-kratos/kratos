package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_VideoAudit(t *testing.T) {
	var (
		c  = context.TODO()
		vp = &archive.VideoParam{}
	)
	attrs := make(map[uint]int32, 5)
	attrs[archive.AttrBitNoRank] = 1
	attrs[archive.AttrBitNoDynamic] = 1
	attrs[archive.AttrBitNoSearch] = 1
	attrs[archive.AttrBitNoRecommend] = 1
	attrs[archive.AttrBitOverseaLock] = 1
	Convey("VideoAudit", t, WithService(func(s *Service) {
		err := svr.VideoAudit(c, vp, attrs)
		So(err, ShouldBeNil)
	}))
}

func TestService_UpVideo(t *testing.T) {
	var (
		c  = context.TODO()
		vp = &archive.VideoParam{}
	)

	Convey("UpVideo", t, WithService(func(s *Service) {
		err := svr.UpVideo(c, vp)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_UpWebLink(t *testing.T) {
	var (
		c  = context.TODO()
		vp = &archive.VideoParam{}
	)

	Convey("UpWebLink", t, WithService(func(s *Service) {
		err := svr.UpWebLink(c, vp)
		So(err, ShouldBeNil)
	}))
}

func TestService_DelVideo(t *testing.T) {
	var (
		c  = context.TODO()
		vp = &archive.VideoParam{}
	)

	Convey("DelVideo", t, WithService(func(s *Service) {
		err := svr.DelVideo(c, vp)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_ChangeIndex(t *testing.T) {
	var (
		c  = context.TODO()
		vp = &archive.IndexParam{}
	)

	Convey("ChangeIndex", t, WithService(func(s *Service) {
		err := svr.ChangeIndex(c, vp)
		So(err, ShouldNotBeNil)
	}))
}
