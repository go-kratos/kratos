package version

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/version"
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

func Test_VersionMap(t *testing.T) {
	var (
		c        = context.Background()
		err      error
		versions = make(map[string][]*version.Version)
	)
	Convey("versionMap", t, WithService(func(s *Service) {
		versions, err = s.versionMap(c)
		So(err, ShouldBeNil)
		So(versions, ShouldNotBeNil)
	}))
}
