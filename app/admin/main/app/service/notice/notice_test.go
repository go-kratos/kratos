package notice

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/notice"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestNotices(t *testing.T) {
	Convey("get Notice", t, WithService(func(s *Service) {
		res, err := s.Notices(context.TODO())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestNoticeByID(t *testing.T) {
	Convey("get TestNoticeByID", t, WithService(func(s *Service) {
		res, err := s.NoticeByID(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestInsert(t *testing.T) {
	Convey("insert notice", t, WithService(func(s *Service) {
		a := &notice.Param{
			Plat:    2,
			Title:   "苏轼",
			Content: "不思量，自难忘",
			URL:     "http://www.bilibili.com",
			Area:    "中国台湾",
			Type:    1,
		}
		err := s.Insert(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdateNotice(t *testing.T) {
	Convey("update notice", t, WithService(func(s *Service) {
		a := &notice.Param{
			ID:      19,
			Plat:    3,
			Title:   "白居易",
			Content: "相见时难别亦难",
			URL:     "http://www.bilibili.com",
			Area:    "中国香港",
			Type:    2,
		}
		err := s.UpdateNotice(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdateBuild(t *testing.T) {
	Convey("UpdateBuild notice", t, WithService(func(s *Service) {
		a := &notice.Param{
			ID:         18,
			Build:      2,
			Conditions: "eq",
		}
		err := s.UpdateBuild(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdateState(t *testing.T) {
	Convey("UpdateRelease notice", t, WithService(func(s *Service) {
		a := &notice.Param{
			ID:    17,
			State: 0,
		}
		err := s.UpdateState(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}
