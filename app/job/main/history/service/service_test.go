package service

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/history/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	dir, _ := filepath.Abs("../cmd/history-job-test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}

func Test_Ping(t *testing.T) {
	Convey("Test_Ping", t, func() {
		s.Ping()
	})
}
