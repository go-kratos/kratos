package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/activity/conf"
	l "go-common/app/job/main/activity/model/like"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/activity-job-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	if svf == nil {
		svf = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svf)
	}
}
func TestService_LikeArc(t *testing.T) {
	Convey("test archive view", t, WithService(func(s *Service) {
		sub := &l.Subject{
			ID: 1,
		}
		res, err := s.likeArc(context.Background(), sub)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
