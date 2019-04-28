package feedback

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
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

func Test_State(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	Convey("AppealList", t, WithService(func(s *Service) {
		_, err := s.Tags(c, mid, "")
		So(err, ShouldNotBeNil)
	}))
	Convey("Detail", t, WithService(func(s *Service) {
		_, err := s.Detail(c, mid, 23333, "")
		So(err, ShouldNotBeNil)
	}))
	Convey("Feedbacks", t, WithService(func(s *Service) {
		_, count, err := s.Feedbacks(c, mid, 10, 10, int64(2), "1", "2", "3", "pc", "")
		So(err, ShouldNotBeNil)
		So(count, ShouldBeZeroValue)
	}))
}
