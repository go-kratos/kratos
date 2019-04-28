package service

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/infra/databus/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/databus-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestArchive(t *testing.T) {
	convey.Convey("Archive", t, func() {
		a, ok := s.AuthApp("databus_test_group")
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(a, convey.ShouldNotBeNil)
		convey.Printf("%+v", a)
	})
}
