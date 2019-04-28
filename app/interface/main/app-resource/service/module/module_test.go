package module

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	"go-common/app/interface/main/app-resource/model/module"

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

func TestList(t *testing.T) {
	Convey("get list data", t, WithService(func(s *Service) {
		res := s.List(context.TODO(), "iphone", "phone", "ios", "resourcefile", "1", 4500, 0, 0, 0, 0, []*module.Versions{}, time.Now())
		result, _ := json.Marshal(res)
		fmt.Printf("test list (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
	}))
}

func TestResource(t *testing.T) {
	Convey("get Resource data", t, WithService(func(s *Service) {
		res, err := s.Resource(context.TODO(), "iphone", "phone", "ios", "resourcefile", "我的测试包222", "1", 21, 3500, 0, 0, 0, 0, time.Now())
		result, _ := json.Marshal(res)
		fmt.Printf("test Resource (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
