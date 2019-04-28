package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/push/conf"
	pushmdl "go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/push-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func Test_Setting(t *testing.T) {
	Convey("setting", t, WithService(func(s *Service) {
		var (
			c   = context.Background()
			mid = int64(91221505)
		)
		err := s.SetSetting(c, mid, pushmdl.UserSettingArchive, pushmdl.SwitchOff)
		So(err, ShouldBeNil)

		setting, err := s.Setting(c, mid)
		So(err, ShouldBeNil)
		st := make(map[int]int, len(pushmdl.Settings))
		for k, v := range pushmdl.Settings {
			st[k] = v
		}
		st[pushmdl.UserSettingArchive] = pushmdl.SwitchOff
		So(setting, ShouldResemble, st)
	}))

	Convey("get default setting", t, WithService(func(s *Service) {
		setting, err := s.Setting(context.TODO(), 8888888888888)
		t.Logf("setting(%+v)", pushmdl.Settings)
		So(err, ShouldBeNil)
		So(setting, ShouldResemble, pushmdl.Settings)
	}))
}

func Benchmark_Callback(b *testing.B) {
	Convey("callback", b, WithService(func(s *Service) {
		// for n := 0; n < b.N; n++ {
		// 	s.CallbackClick(context.TODO(), &pushmdl.Callback{
		// 		Type: pushmdl.CallbackTypeClick,
		// 	})
		// }
	}))
}

func TestServicever2build(t *testing.T) {
	version := "5.7.1(5730)"
	res := ver2build(version, pushmdl.PlatformIPhone)
	if res != 5730 {
		t.FailNow()
	}
	version = "5.7.1"
	res = ver2build(version, pushmdl.PlatformIPhone)
	if res != 5730 {
		t.FailNow()
	}
	version = "5.14.0"
	res = ver2build(version, pushmdl.PlatformAndroid)
	if res != 514000 {
		t.FailNow()
	}
	version = "5.14.0-preview"
	res = ver2build(version, pushmdl.PlatformAndroid)
	if res != 514000 {
		t.FailNow()
	}
}
