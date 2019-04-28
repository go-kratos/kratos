package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/stat/conf"
	"go-common/app/service/main/archive/api"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/stat-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		s.Ping(context.TODO())
	})
}

func Test_UpdateCache(t *testing.T) {
	Convey("updateCache", t, func() {
		err := s.updateCache(&api.Stat{
			Aid:  1,
			Coin: 22,
			View: 200,
		})
		So(err, ShouldBeNil)
	})
}
