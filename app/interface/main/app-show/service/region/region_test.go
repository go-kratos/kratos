package region

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-show/conf"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestRegions(t *testing.T) {
	Convey("get Regions data", t, WithService(func(s *Service) {
		res, ver, err := s.Regions(context.TODO(), 0, 11111, "", "android", "", _initlanguage)
		So(res, ShouldNotBeEmpty)
		So(ver, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestRegionsList(t *testing.T) {
	Convey("get RegionsList data", t, WithService(func(s *Service) {
		res, ver, err := s.RegionsList(context.TODO(), 0, 11111, "", "android", "", _initlanguage, "region")
		So(res, ShouldNotBeEmpty)
		So(ver, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
