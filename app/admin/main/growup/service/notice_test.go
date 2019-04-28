package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InsertNotice(t *testing.T) {
	title := "公告标题"
	typ, platform, status := 0, 1, 0
	link := "www.bilibili.com"
	Convey("admins", t, WithService(func(s *Service) {
		err := s.InsertNotice(context.Background(), title, typ, platform, link, status)
		So(err, ShouldBeNil)
	}))
}

func Test_Notices(t *testing.T) {
	typ, platform, status := 0, 1, 0
	from, limit := 0, 1000
	Convey("admins", t, WithService(func(s *Service) {
		_, res, err := s.Notices(context.Background(), typ, status, platform, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func Test_UpdateNotice(t *testing.T) {
	title := "公告标题"
	typ, platform, status, id := 0, 1, 0, int64(0)
	link := "www.bilibili.com"
	Convey("admins", t, WithService(func(s *Service) {
		err := s.UpdateNotice(context.Background(), typ, platform, title, link, id, status)
		So(err, ShouldBeNil)
	}))
}
