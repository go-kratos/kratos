package assist

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/assist"

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

func ctx() context.Context {
	return context.Background()
}

func Test_Revoc(t *testing.T) {
	var (
		c         = ctx()
		ck        = ""
		ip        = "127.0.0.1"
		assistLog = &assist.AssistLog{}
	)
	Convey("test_revoc", t, func() {
		err := s.revoc(c, assistLog, ck, ip)
		So(err, ShouldNotBeNil)
	})
}

func Test_Live(t *testing.T) {
	var (
		c         = ctx()
		ip        = "127.0.0.1"
		ck        = ""
		ak        = ""
		MID       = int64(27515256)
		assistMid = int64(27515257)
	)
	Convey("test_LiveStatus", t, func() {
		open, err := s.LiveStatus(c, MID, ip)
		So(err, ShouldBeNil)
		So(open, ShouldNotBeNil)
	})
	Convey("test_liveAddAssist", t, func() {
		err := s.liveAddAssist(c, MID, assistMid, ak, ck, ip)
		So(err, ShouldBeNil)
	})
}

func Test_Assist(t *testing.T) {
	var (
		c   = ctx()
		err error
		ip  = "127.0.0.1"
		MID = int64(27515256)
		ass []*assist.Assist
	)
	Convey("test_Assists", t, WithService(func(s *Service) {
		ass, err = s.Assists(c, MID, ip)
		So(err, ShouldBeNil)
		So(ass, ShouldNotBeNil)
	}))
}
