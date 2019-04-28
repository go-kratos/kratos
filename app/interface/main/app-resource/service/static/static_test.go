package static

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-common/app/interface/main/app-resource/conf"
	"path/filepath"
	"testing"
	"time"

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

// go test -conf="../../cmd/app-resource-test.toml"  -v -test.run TestStatic
func TestStatic(t *testing.T) {
	Convey("get static data", t, WithService(func(s *Service) {
		res, ver, err := s.Static(1, 22222, "", time.Now())
		result, err := json.Marshal(res)
		fmt.Printf("test static (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
		So(ver, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
