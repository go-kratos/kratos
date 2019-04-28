package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/tv/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
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

func TestService_SearchTypes(t *testing.T) {
	Convey("TestService_SearchTypes", t, WithService(func(s *Service) {
		cont, err := s.SearchTypes()
		So(err, ShouldBeNil)
		So(len(cont), ShouldBeGreaterThan, 0)
		for _, v := range cont {
			data, _ := json.Marshal(v)
			fmt.Println(string(data))
		}
	}))
}
