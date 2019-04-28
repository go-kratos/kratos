package service

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/block/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c context.Context
)

func TestMain(m *testing.M) {
	defer os.Exit(0)
	flag.Set("conf", "../cmd/block-job-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New()
	defer s.Close()

	m.Run()
}

func TestService(t *testing.T) {
	Convey("", t, func() {
		s.Ping(c)
		s.Close()
	})
}
