package domain

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

func TestDomain(t *testing.T) {
	Convey("get Domain data", t, WithService(func(s *Service) {
		res := s.Domain()
		result, _ := json.Marshal(res)
		fmt.Printf("test Domain (%v) \n", string(result))
		So(res, ShouldNotBeEmpty)
	}))
}
