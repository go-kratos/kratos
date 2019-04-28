package daily

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"

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
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestDaily(t *testing.T) {
	Convey("get Daily data", t, WithService(func(s *Service) {
		res := s.Daily(context.TODO(), model.PlatIPhone, 100000, 4, 1, 20)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestColumnList(t *testing.T) {
	Convey("get ColumnList data", t, WithService(func(s *Service) {
		res := s.ColumnList(model.PlatIPhone, 100000, 4)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestCategory(t *testing.T) {
	Convey("get Category data", t, WithService(func(s *Service) {
		res := s.Category(model.PlatIPhone, 100000, 4, 4, 1, 20)
		So(res, ShouldNotBeEmpty)
	}))
}
