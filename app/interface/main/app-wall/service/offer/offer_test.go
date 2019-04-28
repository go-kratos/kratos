package offer

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func TestTelecomPay(t *testing.T) {
	Convey("TelecomPay", t, WithService(func(s *Service) {
		err := s.Click(context.TODO(), 1, "", "", "", "", time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestANClick(t *testing.T) {
	Convey("ANClick", t, WithService(func(s *Service) {
		err := s.ANClick(context.TODO(), "", "", "", "", "", 1, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestANActive(t *testing.T) {
	Convey("ANActive", t, WithService(func(s *Service) {
		err := s.ANActive(context.TODO(), "", "", "", time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestActive(t *testing.T) {
	Convey("Active", t, WithService(func(s *Service) {
		err := s.Active(context.TODO(), 1, "", "", "", "", "", time.Now())
		So(err, ShouldBeNil)
	}))
}
