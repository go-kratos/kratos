package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/coin/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/coin-job-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic("conf.Init() error")
	}
	s = New(conf.Conf)
	time.Sleep(time.Second * 1)
}

func TestService(t *testing.T) {
	Convey("mock", t, func() {
		s.Close()
		s.Redo(0)
		s.Ping(context.TODO())
		s.Wait()
	})
}
