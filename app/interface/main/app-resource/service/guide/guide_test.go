package guide

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

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

func TestGuide(t *testing.T) {
	Convey("get guide data", t, WithService(func(s *Service) {
		res := s.Interest("iphone", "", time.Now())
		result, _ := json.Marshal(res)
		fmt.Printf("test guide (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_Guide2(t *testing.T) {
	Convey("Guide2", t, WithService(func(s *Service) {
		res := s.Interest2("ssss11", "ssss11", "ssss11", "ssss11", 1, time.Now())
		result, _ := json.Marshal(res)
		Printf("test guide (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
	}))
}
