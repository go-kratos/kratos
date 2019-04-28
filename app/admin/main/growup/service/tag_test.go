package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/growup-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		// Reset(func() { CleanCache() })
		f(srv)
	}
}

func Test_AddTagInfo(t *testing.T) {
	var (
		c   = context.Background()
		tag = &model.TagInfo{
			Tag:      "标签1",
			Category: 6,
			Business: 1,
		}
		creator = "creator"
	)
	Convey("admins", t, WithService(func(s *Service) {
		s.dao.Exec(c, "DELETE FROM tag_info WHERE tag = '标签1'")
		err := s.AddTagInfo(c, tag, creator)
		So(err, ShouldBeNil)
	}))
}

func Test_ModeTagState(t *testing.T) {
	var (
		c         = context.Background()
		tagID     = 1
		isDeleted = 1
	)
	Convey("admins", t, WithService(func(s *Service) {
		s.dao.Exec(c, "DELETE FROM tag_info WHERE id = 1")
		s.dao.Exec(c, "insert into tag_info(id, tag) values(1, 'test')")
		err := s.ModeTagState(c, tagID, isDeleted)
		So(err, ShouldBeNil)
	}))
}

func Test_AddTagUps(t *testing.T) {
	var (
		c        = context.Background()
		tagID    = 1
		mids     = []int64{1011, 1022}
		isCommon = 1
	)

	Convey("admins", t, WithService(func(s *Service) {
		s.dao.Exec(c, "DELETE FROM tag_up_info WHERE mid in (1011, 1022)")
		err := s.AddTagUps(c, tagID, mids, isCommon)
		So(err, ShouldBeNil)
	}))
}

func Test_ReleaseUp(t *testing.T) {
	var (
		c     = context.Background()
		tagID = 1
		mid   = int64(1011)
	)
	Convey("admins", t, WithService(func(s *Service) {
		s.dao.Exec(c, "UPDATE tag_info SET is_common = 0 WHERE id = 1")
		err := s.ReleaseUp(c, tagID, mid)
		So(err, ShouldBeNil)
	}))
}

func Test_QueryTagInfo(t *testing.T) {
	var (
		c          = context.Background()
		startTime  = int64(1515945600)
		endTime    = int64(1516809600)
		categories = []int64{21, 167}
		business   = []int64{1, 2}
		tag        = "标签1"
		from       = 0
		limit      = 20
		sort       = "-ctime"
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.QueryTagInfo(c, startTime, endTime, categories, business, tag, 0, from, limit, sort)
		So(err, ShouldBeNil)
	}))
}

func Test_ListUps(t *testing.T) {
	var (
		c     = context.Background()
		tagID = 1
		mid   = int64(1011)
		from  = 0
		limit = 20
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.ListUps(c, tagID, mid, from, limit)
		So(err, ShouldBeNil)
	}))
}

func Test_ListAvs(t *testing.T) {
	var (
		c     = context.Background()
		tagID = 1
		avid  = int64(2011)
		from  = 0
		limit = 20
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, _, err := s.ListAvs(c, tagID, from, limit, avid)
		So(err, ShouldBeNil)
	}))
}

func Test_TagDetails(t *testing.T) {
	var (
		c     = context.Background()
		tagID = 1
		from  = 0
		limit = 20
	)
	Convey("admins", t, WithService(func(s *Service) {
		_, res, err := s.TagDetails(c, tagID, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		t.Logf("admins len(%d)", len(res))
	}))
}
