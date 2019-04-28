package service

import (
	"flag"
	"testing"

	"go-common/app/job/main/relation/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		panic(err)
	}

	s = New(conf.Conf)
}

func TestService_Close(t *testing.T) {
	Convey("Close", t, func() {
		So(s.Close(), ShouldBeNil)
	})
}
