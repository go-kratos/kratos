package service

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/credit/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestService_loadCase(t *testing.T) {
	Convey("should return err be nil", t, func() {
		s.loadConf()
		s.loadCase()
	})
}
