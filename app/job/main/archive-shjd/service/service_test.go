package service

import (
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/archive-shjd/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/archive-job-kisjd-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		s.Ping()
	})
}

func Test_DelteVideoCache(t *testing.T) {
	Convey("DelteVideoCache", t, func() {
		s.DelteVideoCache(1, 1)
	})
}

func Test_UpdateVideoCache(t *testing.T) {
	Convey("UpdateVideoCache", t, func() {
		s.UpdateVideoCache(1, 1)
	})
}

func Test_Close(t *testing.T) {
	Convey("Close", t, func() {
		s.Close()
	})
}
