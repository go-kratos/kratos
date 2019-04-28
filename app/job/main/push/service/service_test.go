package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/push/conf"
	pushmdl "go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var srv *Service

func init() {
	dir, _ := filepath.Abs("../cmd/push-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func Test_Ping(t *testing.T) {
	Convey("ping", t, WithService(func(s *Service) {
		err := s.Ping(context.TODO())
		So(err, ShouldBeNil)
	}))
}

func Test_TxCond(t *testing.T) {
	Convey("query conditon by tx", t, WithService(func(s *Service) {
		cond, err := s.txCond(pushmdl.DpCondStatusPrepared, pushmdl.DpCondStatusSubmitting)
		So(err, ShouldBeNil)
		t.Logf("cond(%+v)", cond)
	}))
}
