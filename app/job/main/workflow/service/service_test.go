package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/workflow/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func WithService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../goconvey.toml")
		flag.Set("conf", dir)
		conf.Init()
		s := New(conf.Conf)
		f(s)
	}
}

func Test_queueproc(t *testing.T) {
	var (
		c        = context.TODO()
		dealType = 1
	)
	Convey("queueproc", t, WithService(func(s *Service) {
		s.queueproc(c, dealType)
		So(nil, ShouldBeNil)
	}))
}

func Test_taskExpireproc(t *testing.T) {
	var (
		c        = context.TODO()
		dealType = 1
	)
	Convey("taskExpireproc", t, WithService(func(s *Service) {
		s.taskExpireproc(c, dealType)
		So(nil, ShouldBeNil)
	}))
}
