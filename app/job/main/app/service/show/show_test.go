package show

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app/conf"

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
	dir, _ := filepath.Abs("../../cmd/app-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestPing(t *testing.T) {
	Convey("get Ping data", t, WithService(func(s *Service) {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	}))
}
