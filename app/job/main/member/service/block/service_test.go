package block

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/member/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c context.Context
)

func TestMain(m *testing.M) {
	defer os.Exit(0)
	flag.Set("conf", "../cmd/member-job-dev.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	m.Run()
}

func TestService(t *testing.T) {
	Convey("", t, func() {
		s.Ping(c)
		s.Close()
	})
}
