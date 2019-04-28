package bfs

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"

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
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestClientUpCover(t *testing.T) {
	Convey("pull ClientUpCover", t, WithService(func(s *Service) {
		res, err := s.ClientUpCover(context.TODO(), "image/jpeg", nil)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
