package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/service/main/archive/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/archive-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_AllTypes(t *testing.T) {
	Convey("AllTypes", t, func() {
		types := s.AllTypes(context.TODO())
		So(types, ShouldNotBeNil)
	})
}

func Test_CacheUpdate(t *testing.T) {
	Convey("CacheUpdate", t, func() {
		err := s.CacheUpdate(context.TODO(), 1, "update", 1)
		So(err, ShouldBeNil)
	})
}

func Test_FieldCacheUpdate(t *testing.T) {
	Convey("FieldCacheUpdate", t, func() {
		err := s.FieldCacheUpdate(context.TODO(), 1, 1, 2)
		So(err, ShouldBeNil)
	})
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}
