package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/tag/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	dir, _ := filepath.Abs("../cmd/tag-job-test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}

func TestPing(t *testing.T) {
	Convey("test ping", t, func() {
		s.Ping(context.TODO())
	})
}
