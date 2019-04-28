package service

import (
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/tv/conf"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/tv-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(srv)
	}
}

func Test_existArcTypes(t *testing.T) {
	Convey("existArcTs", t, WithService(func(s *Service) {
		exist, err := s.existArcTps(true)
		fmt.Println(exist)
		fmt.Println(err)
		So(err, ShouldBeNil)
		exist, err = s.existArcTps(false)
		fmt.Println(exist)
		fmt.Println(err)
		So(err, ShouldBeNil)
	}))
}

func Test_Wait(t *testing.T) {
	Convey("wait all closed", t, WithService(func(s *Service) {
		s.Wait()
	}))
}

func TestService_ResExist(t *testing.T) {
	Convey("res exist", t, WithService(func(s *Service) {
		fmt.Println(s.resExist(255, 1))
		fmt.Println(s.recomExist(2))
	}))
}
