package splash

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model"

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
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestDisplay(t *testing.T) {
	Convey("get Display data", t, WithService(func(s *Service) {
		res, _, err := s.Display(context.TODO(), model.PlatIPhone, 1, 1, 1, "", "", time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestBirthday(t *testing.T) {
	Convey("get Birthday data", t, WithService(func(s *Service) {
		res, err := s.Birthday(context.TODO(), model.PlatIPhone, 1, 1, "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
