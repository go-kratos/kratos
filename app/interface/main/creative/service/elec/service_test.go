package elec

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_State(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	Convey("ArchiveState", t, WithService(func(s *Service) {
		_, err := s.ArchiveState(c, 2333, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("UserState", t, WithService(func(s *Service) {
		_, err := s.UserState(c, mid, "", "", "")
		So(err, ShouldBeNil)
	}))
	Convey("UserInfo", t, WithService(func(s *Service) {
		_, err := s.UserInfo(c, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("Notify", t, WithService(func(s *Service) {
		_, err := s.Notify(c, "")
		So(err, ShouldBeNil)
	}))
	Convey("Status", t, WithService(func(s *Service) {
		_, err := s.Status(c, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("UpStatus", t, WithService(func(s *Service) {
		err := s.UpStatus(c, mid, 1, "")
		So(err, ShouldBeNil)
	}))
	Convey("RecentRank", t, WithService(func(s *Service) {
		_, err := s.RecentRank(c, mid, 10, "")
		So(err, ShouldBeNil)
	}))
	Convey("CurrentRank", t, WithService(func(s *Service) {
		_, err := s.CurrentRank(c, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("TotalRank", t, WithService(func(s *Service) {
		_, err := s.TotalRank(c, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("Balance", t, WithService(func(s *Service) {
		_, err := s.Balance(c, mid, "")
		So(err, ShouldBeNil)
	}))
	Convey("RemarkDetail", t, WithService(func(s *Service) {
		_, err := s.RemarkDetail(c, mid, 233, "")
		So(err, ShouldBeNil)
	}))
}
