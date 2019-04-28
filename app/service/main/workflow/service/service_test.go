package service

import (
	"context"
	"flag"
	"testing"

	"go-common/app/service/main/workflow/conf"

	"github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		panic(err)
	}

	s = New(conf.Conf)
}

func TestPing(t *testing.T) {
	convey.Convey("Ping", t, func() {
		err := s.Ping(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}
