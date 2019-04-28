package danmu

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
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

func Test_AdvDmPurchaseList(t *testing.T) {
	var (
		c   = context.Background()
		MID = int64(27515256)
	)
	Convey("AdvDmPurchaseList", t, WithService(func(s *Service) {
		_, err := s.AdvDmPurchaseList(c, MID, "")
		So(err, ShouldNotBeNil)
	}))
}

func Test_Distri(t *testing.T) {
	var (
		c   = context.Background()
		MID = int64(27515256)
	)
	Convey("Distri", t, WithService(func(s *Service) {
		_, err := s.Distri(c, MID, 2333, "")
		So(err, ShouldNotBeNil)
	}))
}

func Test_DmProtectList(t *testing.T) {
	var (
		c   = context.Background()
		MID = int64(27515256)
	)
	Convey("DmProtectList", t, WithService(func(s *Service) {
		_, err := s.DmProtectList(c, MID, 1, "2333", "time", "")
		So(err, ShouldNotBeNil)
	}))
}

func Test_DmReportList(t *testing.T) {
	var (
		c   = context.Background()
		MID = int64(27515256)
	)
	Convey("DmReportList", t, WithService(func(s *Service) {
		_, err := s.DmReportList(c, MID, 1, 10, "2333", "")
		So(err, ShouldNotBeNil)
	}))
}

func Test_Recent(t *testing.T) {
	var (
		c   = context.Background()
		MID = int64(27515256)
	)
	Convey("Recent", t, WithService(func(s *Service) {
		_, err := s.Recent(c, MID, 1, 10, "")
		So(err, ShouldNotBeNil)
	}))
}
