package version

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

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
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestVersion(t *testing.T) {
	Convey("get Version data", t, WithService(func(s *Service) {
		res, err := s.Version(1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestVersionUpdate(t *testing.T) {
	Convey("get VersionUpdate data", t, WithService(func(s *Service) {
		res, err := s.VersionUpdate(1, 1, "", "", "", "", "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestVersionSo(t *testing.T) {
	Convey("get VersionSo data", t, WithService(func(s *Service) {
		res, err := s.VersionSo(1, 1, 1, "", "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestVersionRn(t *testing.T) {
	Convey("get VersionRn data", t, WithService(func(s *Service) {
		res, err := s.VersionRn("", "", "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
