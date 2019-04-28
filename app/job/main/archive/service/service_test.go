package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/archive/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/archive-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_loadType(t *testing.T) {
	Convey("loadType", t, func() {
		s.loadType()
	})
}

func Test_PopFail(t *testing.T) {
	Convey("PopFail", t, func() {
		s.PopFail(context.TODO())
	})
}

func Test_TranResult(t *testing.T) {
	Convey("tranResult", t, func() {
		_, _, _, err := s.tranResult(context.TODO(), 10098500)
		So(err, ShouldBeNil)
	})
}
