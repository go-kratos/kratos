package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/videoup/conf"
	"go-common/library/queue/databus/report"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/videoup-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	svr = New(conf.Conf)
	time.Sleep(time.Second)
	report.InitManager(conf.Conf.ManagerReport)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}

func TestService_Ping(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("Ping", t, WithService(func(s *Service) {
		err := svr.Ping(c)
		So(err, ShouldBeNil)
	}))
}

func TestService_PGCWhite(t *testing.T) {
	Convey("PGCWhite", t, WithService(func(s *Service) {
		data := svr.PGCWhite(1)
		So(data, ShouldBeFalse)
	}))
}

func TestService_getallupgroups(t *testing.T) {
	Convey("getallupgroups", t, WithService(func(s *Service) {
		gs := s.getAllUPGroups(27515615)
		t.Logf("UPGroups(%+v)\r\n", gs)
		So(len(gs), ShouldBeGreaterThan, 0)
	}))
}
