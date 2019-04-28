package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/favorite/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/favorite-job-test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}

func Test_archiveRPC(t *testing.T) {
	Convey("archiveRPC", t, func() {
		var (
			aid int64 = 123
		)
		res, err := s.archiveRPC(context.TODO(), aid)
		t.Logf("res:%+v", res)
		t.Logf("err:%v", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
func Test_ArcsRPC(t *testing.T) {
	Convey("ArcsRPC", t, func() {
		var (
			aids = []int64{123, 456}
		)
		res, err := s.ArcsRPC(context.TODO(), aids)
		t.Logf("res:%+v", res)
		t.Logf("err:%v", err)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
