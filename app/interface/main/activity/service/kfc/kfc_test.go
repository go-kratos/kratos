package kfc

import (
	"testing"
	"time"

	"flag"
	"go-common/app/interface/main/activity/conf"
	"path/filepath"

	"context"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/activity-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	if svf == nil {
		svf = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svf)
	}
}

func TestService_KfcInfo(t *testing.T) {
	Convey("test fmt start and end", t, WithService(func(s *Service) {
		id := int64(30)
		mid := int64(16299551)
		start, err := s.KfcInfo(context.Background(), id, mid)
		So(err, ShouldBeNil)
		Println(start)
	}))
}

func TestService_KfcUse(t *testing.T) {
	Convey("test fmt start and end", t, WithService(func(s *Service) {
		code := "535487458740"
		start, err := s.KfcUse(context.Background(), code)
		So(err, ShouldBeNil)
		Println(start)
	}))
}

func TestService_DeliverKfc(t *testing.T) {
	Convey("test fmt start and end", t, WithService(func(s *Service) {
		id := int64(1)
		mid := int64(2089810)
		err := s.DeliverKfc(context.Background(), id, mid)
		So(err, ShouldBeNil)
	}))
}

func TestService_kfcRecall(t *testing.T) {
	Convey("test fmt start and end", t, WithService(func(s *Service) {
		id := int64(30)
		uid, err := s.kfcRecall(context.Background(), id)
		So(err, ShouldBeNil)
		fmt.Print(uid)
	}))
}
