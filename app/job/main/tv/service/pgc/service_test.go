package pgc

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/tv/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func TestService_SearchSug(t *testing.T) {
	Convey("TestService_SearchSug", t, WithService(func(s *Service) {
		s.searchSug()
	}))
}
