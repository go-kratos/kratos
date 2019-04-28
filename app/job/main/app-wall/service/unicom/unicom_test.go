package unicom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app-wall/conf"
	"go-common/app/job/main/app-wall/model/unicom"

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
	dir, _ := filepath.Abs("../../cmd/app-wall-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAddUserPackLogproc(t *testing.T) {
	Convey("Unicom addUserPackLogproc", t, WithService(func(s *Service) {
		s.addUserPackLogproc()
	}))
}

func TestAddUserIntegralLogproc(t *testing.T) {
	Convey("Unicom addUserIntegralLogproc", t, WithService(func(s *Service) {
		s.addUserIntegralLogproc(1)
	}))
}

func TestAddUserIntegralLog(t *testing.T) {
	Convey("Unicom addUserIntegralLog", t, WithService(func(s *Service) {
		s.addUserIntegralLog(&unicom.UserIntegralLog{})
	}))
}

func TestLoadUnicomIP(t *testing.T) {
	Convey("Unicom loadUnicomIP", t, WithService(func(s *Service) {
		s.loadUnicomIP(context.TODO())
	}))
}

func TestLoadUnicomIPOrder(t *testing.T) {
	Convey("Unicom loadUnicomIPOrder", t, WithService(func(s *Service) {
		s.loadUnicomIPOrder(time.Now())
	}))
}
